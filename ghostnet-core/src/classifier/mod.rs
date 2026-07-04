use crate::common::{Packet, TrafficClass};
use std::collections::HashMap;
use std::sync::{Arc, RwLock};

pub struct Policy {
    pub merge_threshold: usize,
    pub split_threshold: usize,
}

pub struct Classifier {
    policy: Arc<RwLock<Policy>>,
    queue_volumes: Arc<RwLock<HashMap<TrafficClass, usize>>>,
}

impl Classifier {
    pub fn new() -> Self {
        Self {
            policy: Arc::new(RwLock::new(Policy {
                merge_threshold: 10,
                split_threshold: 100,
            })),
            queue_volumes: Arc::new(RwLock::new(HashMap::new())),
        }
    }

    pub fn classify(&self, p: &Packet) -> TrafficClass {
        let size = p.size;
        let payload = String::from_utf8_lossy(&p.payload);

        let class = if payload.contains("GET") && size < 5000 {
            TrafficClass::SmallWeb
        } else if payload.contains("API") {
            TrafficClass::API
        } else if size > 1_000_000 {
            TrafficClass::Video
        } else if size > 100_000 {
            TrafficClass::Image
        } else if size > 10_000 {
            TrafficClass::FileUpload
        } else if payload.contains("DNS") {
            TrafficClass::DNS
        } else if size < 1000 {
            TrafficClass::Background
        } else {
            TrafficClass::Interactive
        };

        let mut volumes = self.queue_volumes.write().unwrap();
        *volumes.entry(class.clone()).or_insert(0) += 1;

        class
    }

    pub fn update_policy(&self, new_policy: Policy) {
        let mut policy = self.policy.write().unwrap();
        *policy = new_policy;
    }

    pub fn get_queues(&self) -> Vec<TrafficClass> {
        let mut volumes = self.queue_volumes.write().unwrap();
        let policy = self.policy.read().unwrap();

        let mut active_classes = Vec::new();
        let mut to_merge = Vec::new();

        for (class, vol) in volumes.iter() {
            if *vol > policy.split_threshold {
                active_classes.push(class.clone());
            } else if *vol < policy.merge_threshold && *vol > 0 {
                to_merge.push(class.clone());
            } else if *vol > 0 {
                active_classes.push(class.clone());
            }
        }

        let mut merged_vol = 0;
        for class in to_merge {
            merged_vol += volumes.remove(&class).unwrap_or(0);
        }

        if merged_vol > 0 {
            *volumes.entry(TrafficClass::Unknown).or_insert(0) += merged_vol;
        }

        if volumes.get(&TrafficClass::Unknown).copied().unwrap_or(0) > 0 {
            active_classes.push(TrafficClass::Unknown);
        }

        active_classes
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::common::{Packet, Metadata};
    use chrono::Utc;

    #[test]
    fn test_classify() {
        let c = Classifier::new();
        let p = Packet {
            id: "test".to_string(),
            payload: b"GET /index.html HTTP/1.1".to_vec(),
            size: 400,
            traffic_class: TrafficClass::Unknown,
            is_dummy: false,
            metadata: Metadata::default(),
            timestamp: Utc::now(),
        };
        assert_eq!(c.classify(&p), TrafficClass::SmallWeb);
    }
}
