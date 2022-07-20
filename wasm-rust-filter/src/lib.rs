use log::trace;
use proxy_wasm::traits::*;
use proxy_wasm::types::*;

proxy_wasm::main! {{
    proxy_wasm::set_log_level(LogLevel::Trace);
    proxy_wasm::set_root_context(|_| -> Box<dyn RootContext> { Box::new(HttpHeadersRoot) });
}}

struct HttpHeadersRoot;

impl Context for HttpHeadersRoot {}

impl RootContext for HttpHeadersRoot {
    fn get_type(&self) -> Option<ContextType> {
        Some(ContextType::HttpContext)
    }

    fn create_http_context(&self, context_id: u32) -> Option<Box<dyn HttpContext>> {
        Some(Box::new(HttpHeaders { context_id }))
    }
}

struct HttpHeaders {
    context_id: u32,
}

impl Context for HttpHeaders {}

impl HttpContext for HttpHeaders {
    fn on_http_request_headers(&mut self, _: usize, _: bool) -> Action {
        for (name, value) in &self.get_http_request_headers() {
            trace!("#{} -> {}: {}", self.context_id, name, value);
        }

        basics();

        match self.get_http_request_header(":path") {
            Some(path) if path == "/hello" => {
                self.send_http_response(
                    200,
                    vec![("Hello", "World"), ("Powered-By", "proxy-wasm")],
                    Some(b"Hello, World!\n"),
                );
                Action::Pause
            }
            _ => Action::Continue,
        }
    }

    fn on_http_response_headers(&mut self, _: usize, _: bool) -> Action {
        for (name, value) in &self.get_http_response_headers() {
            trace!("#{} <- {}: {}", self.context_id, name, value);
        }
        Action::Continue
    }

    fn on_log(&mut self) {
        trace!("#{} completed.", self.context_id);
    }
}

fn connect() -> redis::Connection {
    //format - host:port
    let redis_host_name = "host.docker.internal:6379";
    let redis_password = "";

    let redis_conn_url = format!("{}://:{}@{}", "redis", redis_password, redis_host_name);
    //println!("{}", redis_conn_url);

    redis::Client::open(redis_conn_url)
        .expect("Invalid connection URL")
        .get_connection()
        .expect("failed to connect to Redis")
}

fn basics() {
    let mut conn = connect();

    let _: () = redis::cmd("SET")
        .arg("foo")
        .arg("bar")
        .query(&mut conn)
        .expect("failed to execute SET for 'foo'");

    let bar: String = redis::cmd("GET")
        .arg("foo")
        .query(&mut conn)
        .expect("failed to execute GET for 'foo'");
}