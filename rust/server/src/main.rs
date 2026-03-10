use std::collections::HashMap;

use axum::extract::{Path, Query};
use axum::response::Json;
use axum::routing::post;
use axum::{Router, routing::get};
use serde_json::{Value, json};

#[tokio::main]
async fn main() {
    let app: Router = Router::new()
        .route("/", get(|| async { "Hello World" }))
        .route("/path/{input}", get(read_path))
        .route("/query", get(read_query))
        .route("/json", post(read_payload));
    let listener = tokio::net::TcpListener::bind("127.0.0.1:8080")
        .await
        .unwrap();
    axum::serve(listener, app).await.unwrap();
}

async fn read_path(Path(input): Path<String>) -> Json<Value> {
    println!("message '{}' sended", input);
    Json(json!({"message": input}))
}

async fn read_query(Query(params): Query<HashMap<String, String>>) {
    for key in params.keys() {
        println!("{}:{}", key, params[key]);
    }
}

async fn read_payload(Json(payload): Json<Value>) {
    match payload.get("message") {
        None => println!("no message..."),
        Some(val) => {
            println!("{}", val.as_str().unwrap())
        }
    }
}
