use rand::Rng;
use std::collections::HashMap;
use std::sync::{Arc, Mutex};

#[derive(Debug, Clone, PartialEq)]
pub enum RelayRole {
    Entry,
    Middle,
    Exit,
    Bridge,
    Directory,
}

#[derive(Debug, Clone)]
pub struct Relay {
    pub id: String,
    pub ip: String,
    pub country: String,
    pub role: RelayRole,
    pub reputation: f64,
    pub load: f64,
}

pub struct Engine {
    relays: Arc<Mutex<HashMap<String, Relay>>>,
}

impl Engine {
    pub fn new() -> Self {
        let e = Self {
            relays: Arc::new(Mutex::new(HashMap::new())),
        };
        e.seed_relays();
        e
    }

    fn seed_relays(&self) {
        let countries = vec!["Germany", "Singapore", "Canada", "USA", "India", "UK", "Japan"];
        let mut relays = self.relays.lock().unwrap();
        for i in 0..50 {
            let role = if i < 10 {
                RelayRole::Entry
            } else if i > 40 {
                RelayRole::Exit
            } else {
                RelayRole::Middle
            };

            let id = format!("relay-{}", i);
            relays.insert(id.clone(), Relay {
                id,
                ip: format!("1.2.3.{}", i),
                country: countries[rand::thread_rng().gen_range(0..countries.len())].to_string(),
                role,
                reputation: 1.0,
                load: rand::thread_rng().gen_range(0.0..1.0),
            });
        }
    }

    pub fn select_route(&self) -> Vec<String> {
        let mut relays = self.relays.lock().unwrap();

        let mut entries = Vec::new();
        let mut middles = Vec::new();
        let mut exits = Vec::new();

        for r in relays.values() {
            if r.load > 0.9 || r.reputation < 0.5 {
                continue;
            }

            match r.role {
                RelayRole::Entry => entries.push(r.id.clone()),
                RelayRole::Middle => middles.push(r.id.clone()),
                RelayRole::Exit => exits.push(r.id.clone()),
                _ => {}
            }
        }

        if entries.is_empty() || middles.len() < 2 || exits.is_empty() {
            return Vec::new();
        }

        let e1 = entries[rand::thread_rng().gen_range(0..entries.len())].clone();
        let m1 = middles[rand::thread_rng().gen_range(0..middles.len())].clone();
        let m2 = middles[rand::thread_rng().gen_range(0..middles.len())].clone();
        let ex = exits[rand::thread_rng().gen_range(0..exits.len())].clone();

        if let Some(r) = relays.get_mut(&e1) { r.load += 0.01; }
        if let Some(r) = relays.get_mut(&m1) { r.load += 0.01; }
        if let Some(r) = relays.get_mut(&m2) { r.load += 0.01; }
        if let Some(r) = relays.get_mut(&ex) { r.load += 0.01; }

        vec![e1, m1, m2, ex]
    }
}
