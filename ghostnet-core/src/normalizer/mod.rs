use crate::common::Packet;

pub const FIXED_PACKET_SIZE: usize = 1024;

pub struct Normalizer;

impl Normalizer {
    pub fn new() -> Self {
        Self
    }

    pub fn normalize(&self, mut p: Packet) -> Vec<Packet> {
        if p.size == FIXED_PACKET_SIZE {
            return vec![p];
        }

        if p.size < FIXED_PACKET_SIZE {
            let mut padded_payload = vec![0u8; FIXED_PACKET_SIZE];
            let copy_len = p.payload.len();
            padded_payload[..copy_len].copy_from_slice(&p.payload);
            p.payload = padded_payload;
            p.size = FIXED_PACKET_SIZE;
            return vec![p];
        }

        let mut fragments = Vec::new();
        for (i, chunk) in p.payload.chunks(FIXED_PACKET_SIZE).enumerate() {
            let mut frag_payload = vec![0u8; FIXED_PACKET_SIZE];
            frag_payload[..chunk.len()].copy_from_slice(chunk);

            fragments.push(Packet {
                id: format!("{}-frag-{}", p.id, i),
                payload: frag_payload,
                size: FIXED_PACKET_SIZE,
                traffic_class: p.traffic_class.clone(),
                is_dummy: p.is_dummy,
                metadata: p.metadata.clone(),
                timestamp: p.timestamp,
            });
        }
        fragments
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::common::{Packet, TrafficClass, Metadata};
    use chrono::Utc;

    #[test]
    fn test_normalize_padding() {
        let n = Normalizer::new();
        let p = Packet {
            id: "test".to_string(),
            payload: vec![1, 2, 3],
            size: 3,
            traffic_class: TrafficClass::SmallWeb,
            is_dummy: false,
            metadata: Metadata::default(),
            timestamp: Utc::now(),
        };
        let packets = n.normalize(p);
        assert_eq!(packets.len(), 1);
        assert_eq!(packets[0].payload.len(), FIXED_PACKET_SIZE);
    }
}
