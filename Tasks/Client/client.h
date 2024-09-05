#ifndef CLIENT_H
#define CLIENT_H

#include <netinet/in.h>

#define SERVER_PORT 8005

int createClient();
int createSocket();
void connectServer(int fd, struct sockaddr_in *serv_addr);
void communicateServer(int fd);
void *handleTerminal(void *arg);
void sendMessage(int fd, const char message[]);

#endif
