#ifndef _SERVER_H
#define _SERVER_H

#include <netinet/in.h>
#include <stdlib.h>
#include <signal.h>

#define SERVER_PORT 8005

extern volatile sig_atomic_t running;

char SERVER_IP[INET_ADDRSTRLEN];

int createServer();
int createSocket();
void bindSocket(int fd, struct sockaddr_in *address);
void listenConnections(int fd);
void acceptConnection(int server_fd, struct sockaddr_in *address);
void *handleClient(void *arg);
void *handleTerminal(void *arg);
void sendMessage(int client_fd, const char message[]);

#endif
