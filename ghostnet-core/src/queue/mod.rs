use crate::common::{Packet, TrafficClass};
use rand::seq::SliceRandom;
use rand::Rng;
use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use tokio::sync::mpsc;
use tokio::time::{self, Duration, Instant};

#[derive(Clone, Debug)]
pub struct Config {
    pub min_batch_size: usize,
    pub max_wait_time: Duration,
    pub random_interval: Duration,
}

pub struct MixQueue {
    config: Config,
    packets: Arc<Mutex<Vec<Packet>>>,
    last_update: Arc<Mutex<Instant>>,
    out: mpsc::Sender<Vec<Packet>>,
}

impl MixQueue {
    pub fn new(config: Config, out: mpsc::Sender<Vec<Packet>>) -> Self {
        Self {
            config,
            packets: Arc::new(Mutex::new(Vec::new())),
            last_update: Arc::new(Mutex::new(Instant::now())),
            out,
        }
    }

    pub fn enqueue(&self, p: Packet) {
        let mut packets = self.packets.lock().unwrap();
        packets.push(p);
    }

    pub async fn run(self) {
        let mut interval = time::interval(self.config.random_interval / 2);
        loop {
            interval.tick().await;

            let jitter = Duration::from_micros(rand::thread_rng().gen_range(0..self.config.random_interval.as_micros() as u64));
            time::sleep(jitter).await;

            let mut packets_to_release = Vec::new();
            {
                let mut packets = self.packets.lock().unwrap();
                let last_update = self.last_update.lock().unwrap();

                let dynamic_min_batch = (self.config.min_batch_size as i32 + rand::thread_rng().gen_range(-1..=2)) as usize;

                if packets.len() >= dynamic_min_batch || (!packets.is_empty() && last_update.elapsed() >= self.config.max_wait_time) {
                    packets.as_mut_slice().shuffle(&mut rand::thread_rng());
                    packets_to_release = packets.drain(..).collect();
                }
            }

            if !packets_to_release.is_empty() {
                {
                    let mut last_update = self.last_update.lock().unwrap();
                    *last_update = Instant::now();
                }
                let _ = self.out.send(packets_to_release).await;
            }
        }
    }
}

pub struct Manager {
    queues: Arc<Mutex<HashMap<TrafficClass, mpsc::Sender<Packet>>>>,
    batch_out: mpsc::Sender<Vec<Packet>>,
}

impl Manager {
    pub fn new(batch_out: mpsc::Sender<Vec<Packet>>) -> Self {
        Self {
            queues: Arc::new(Mutex::new(HashMap::new())),
            batch_out,
        }
    }

    pub fn enqueue(&self, p: Packet) {
        let mut queues = self.queues.lock().unwrap();
        let class = p.traffic_class.clone();

        if let Some(tx) = queues.get(&class) {
            let _ = tx.try_send(p);
        } else {
            let (tx, mut rx) = mpsc::channel(100);
            queues.insert(class.clone(), tx.clone());

            let config = Config {
                min_batch_size: 5,
                max_wait_time: Duration::from_millis(500),
                random_interval: Duration::from_millis(100),
            };

            let batch_out = self.batch_out.clone();
            tokio::spawn(async move {
                let q = MixQueue::new(config, batch_out);
                let q_ref = Arc::new(q);
                let q_inner = q_ref.clone();

                tokio::spawn(async move {
                    q_inner.as_ref().clone_for_run().run().await;
                });

                while let Some(p) = rx.recv().await {
                    q_ref.enqueue(p);
                }
            });
            let _ = tx.try_send(p);
        }
    }
}

impl MixQueue {
    // Helper to bypass the need for Clone if I don't want to derive it for everything
    fn clone_for_run(&self) -> Self {
        Self {
            config: self.config.clone(),
            packets: self.packets.clone(),
            last_update: self.last_update.clone(),
            out: self.out.clone(),
        }
    }
}
