use serde::{Serialize, Deserialize};
use chrono::{DateTime, Utc};

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq, Hash)]
pub enum TrafficClass {
    SmallWeb,
    API,
    Image,
    Video,
    FileUpload,
    DNS,
    Background,
    Interactive,
    Unknown,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum PrivacyMode {
    Off,
    Low,
    Medium,
    High,
    MaximumPrivacy,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Packet {
    pub id: String,
    pub payload: Vec<u8>,
    pub size: usize,
    pub traffic_class: TrafficClass,
    pub is_dummy: bool,
    pub metadata: Metadata,
    pub timestamp: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct Metadata {
    pub source: String,
    pub destination: String,
    pub identity_id: String,
    pub route: Vec<String>,
    pub hop_index: usize,
}
