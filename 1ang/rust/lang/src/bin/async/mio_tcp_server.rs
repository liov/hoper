use std::error::Error;

use mio::net::{TcpListener, TcpStream};
use mio::{Events, Registration, PollOpt, Ready, Poll, Token, Event};
use std::collections::HashMap;
use std::io;
use std::str::from_utf8;
use std::convert::TryFrom;
use std::io::{Write, Read};

// Some tokens to allow us to identify which event is for which socket.
const SERVER: Token = Token(0);
const DATA: &[u8] = b"Hello world!\n";

fn main() -> Result<(), Box<dyn Error>> {
    // Create a poll instance.
    let mut poll = Poll::new()?;
    // Create storage for events.
    let mut events = Events::with_capacity(128);

    // Setup the server socket.
    let addr = "127.0.0.1:9000".parse()?;
    let server = TcpListener::bind(&addr)?;
    // Start listening for incoming connections.
    poll.register(&server, SERVER, Ready::readable(),
                  PollOpt::edge())?;

    // Map of `Token` -> `TcpStream`.
    let mut connections = HashMap::new();
    // Unique token for each incoming connection.
    let mut unique_token = Token(SERVER.0 + 1);

    println!("You can connect to the server using `nc`:");
    println!(" $ nc 127.0.0.1 9000");
    println!("You'll see our welcome message and anything you type we'll be printed here.");
    // Start an event loop.
    loop {
        // Poll Mio for events, blocking until we get an event.
        poll.poll(&mut events, None)?;

        // Process each event.
        for event in events.iter() {
            // We can use the token we previously provided to `register` to
            // determine for which socket the event is.
            match event.token() {
                SERVER => {
                    let (connection, address) = server.accept()?;
                    println!("Accepted connection from: {}", address);

                    let client = next(&mut unique_token);
                    poll.register(
                        &connection,
                        client,
                        Ready::readable() | Ready::writable(),
                        PollOpt::edge(),
                    )?;

                    connections.insert(client, connection);
                }
                token => {
                    // (maybe) received an event for a TCP connection.
                    let done = if let Some(connection) = connections.get_mut(&token) {
                        handle_connection_event(&poll, connection, &event)?
                    } else {
                        // Sporadic events happen.
                        false
                    };
                    if done {
                        connections.remove(&token);
                    }
                }
                // We don't expect any events with tokens other then we provided.
                _ => unreachable!(),
            }
        }
    }
}

fn handle_connection_event(poll: &Poll, connection: &mut TcpStream, event: &Event) -> io::Result<bool> {
    let ready = event.readiness();
    if ready.is_writable() {
        // We can (maybe) write to the connection.
        match connection.write(DATA) {
            // We want to write the entire `DATA` buffer in a single go. If we
            // write less we'll return a short write error (same as
            // `io::Write::write_all` does).
            Ok(n) if n < DATA.len() => return Err(io::ErrorKind::WriteZero.into()),
            Ok(_) => {
                // After we've written something we'll reregister the connection
                // to only respond to readable events.
                poll.reregister(connection, event.token(), Ready::readable(), PollOpt::edge())?
            }
            // Would block "errors" are the OS's way of saying that the
            // connection is not actually ready to perform this I/O operation.
            Err(ref err) if would_block(err) => {}
            // Got interrupted (how rude!), we'll try again.
            Err(ref err) if interrupted(err) => {
                return handle_connection_event(poll, connection, event);
            }
            // Other errors we'll consider fatal.
            Err(err) => return Err(err),
        }
    }

    if ready.is_readable() {
        let mut connection_closed = false;
        let mut received_data = Vec::with_capacity(4096);
        // We can (maybe) read from the connection.
        loop {
            let mut buf = [0; 256];
            match connection.read(&mut buf) {
                Ok(0) => {
                    // Reading 0 bytes means the other side has closed the
                    // connection or is done writing, then so are we.
                    connection_closed = true;
                    break;
                }
                Ok(n) => received_data.extend_from_slice(&buf[..n]),
                // Would block "errors" are the OS's way of saying that the
                // connection is not actually ready to perform this I/O operation.
                Err(ref err) if would_block(err) => break,
                Err(ref err) if interrupted(err) => continue,
                // Other errors we'll consider fatal.
                Err(err) => return Err(err),
            }
        }

        if let Ok(str_buf) = from_utf8(&received_data) {
            println!("Received data: {}", str_buf.trim_end());
        } else {
            println!("Received (none UTF-8) data: {:?}", &received_data);
        }

        if connection_closed {
            println!("Connection closed");
            return Ok(true);
        }
    }

    Ok(false)
}

fn next(current: &mut Token) -> Token {
    let next = current.0;
    current.0 += 1;
    Token(next)
}

fn would_block(err: &io::Error) -> bool {
    err.kind() == io::ErrorKind::WouldBlock
}

fn interrupted(err: &io::Error) -> bool {
    err.kind() == io::ErrorKind::Interrupted
}
