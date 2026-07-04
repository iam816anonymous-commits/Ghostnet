use crate::common::{Packet, PrivacyMode, TrafficClass, Metadata};
use rand::Rng;
use std::sync::{Arc, Mutex};
use tokio::sync::mpsc;
use tokio::time::{self, Duration};
use chrono::Utc;

pub struct Engine {
    mode: Arc<Mutex<PrivacyMode>>,
    out: mpsc::Sender<Packet>,
    active_users: Arc<Mutex<usize>>,
}

impl Engine {
    pub fn new(out: mpsc::Sender<Packet>) -> Self {
        Self {
            mode: Arc::new(Mutex::new(PrivacyMode::Medium)),
            out,
            active_users: Arc::new(Mutex::new(0)),
        }
    }

    pub fn set_mode(&self, mode: PrivacyMode) {
        let mut m = self.mode.lock().unwrap();
        *m = mode;
    }

    pub fn set_active_users(&self, count: usize) {
        let mut c = self.active_users.lock().unwrap();
        *c = count;
    }

    pub async fn run(&self) {
        loop {
            let interval = self.get_interval();
            time::sleep(interval).await;

            let mode = {
                let m = self.mode.lock().unwrap();
                m.clone()
            };

            if let PrivacyMode::Off = mode {
                continue;
            }

            let pkt = self.generate_dummy();
            let _ = self.out.send(pkt).await;
        }
    }

    fn get_interval(&self) -> Duration {
        let mode = self.mode.lock().unwrap();
        let users = self.active_users.lock().unwrap();

        let mut base = match *mode {
            PrivacyMode::Off => Duration::from_secs(3600),
            PrivacyMode::Low => Duration::from_secs(2),
            PrivacyMode::Medium => Duration::from_secs(1),
            PrivacyMode::High => Duration::from_millis(500),
            PrivacyMode::MaximumPrivacy => Duration::from_millis(100),
        };

        if *users > 1000 {
            base *= 2;
        }

        base
    }

    fn generate_dummy(&self) -> Packet {
        let mut payload = vec![0u8; 512];
        rand::thread_rng().fill(&mut payload[..]);

        Packet {
            id: "dummy".to_string(),
            payload,
            size: 512,
            traffic_class: TrafficClass::Background,
            is_dummy: true,
            metadata: Metadata::default(),
            timestamp: Utc::now(),
        }
    }
}
