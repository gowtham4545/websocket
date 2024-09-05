#include <stdio.h>
#include "server.h"
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <pthread.h>

#define BUFFER_SIZE 1024

volatile sig_atomic_t running = 0;

int createServer()
{
    int fd = createSocket();
    struct sockaddr_in addr;
    bindSocket(fd, &addr);
    inet_ntop(AF_INET, &addr.sin_addr, SERVER_IP, sizeof(SERVER_IP));
    listenConnections(fd);
    printf("Server started at %s:%d\n", SERVER_IP, SERVER_PORT);
    running = 1;
    return fd;
}

int createSocket()
{
    int fd = socket(AF_INET, SOCK_STREAM, 0);
    if (fd == 0)
    {
        perror("socket create failed");
        exit(EXIT_FAILURE);
    }

    int opt = 1;
    if (setsockopt(fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt)))
    {
        perror("setsockopt failed");
        close(fd);
        exit(EXIT_FAILURE);
    }

    return fd;
}

void bindSocket(int fd, struct sockaddr_in *address)
{
    address->sin_family = AF_INET;
    address->sin_addr.s_addr = INADDR_ANY;
    address->sin_port = htons(SERVER_PORT);

    if (bind(fd, (struct sockaddr *)address, sizeof(*address)) < 0)
    {
        perror("bind failed");
        close(fd);
        exit(EXIT_FAILURE);
    }
}

void listenConnections(int fd)
{
    if (listen(fd, 3) < 0)
    {
        perror("listen failed");
        close(fd);
        exit(EXIT_FAILURE);
    }
}

void acceptConnection(int server_fd, struct sockaddr_in *address)
{
    int addrlen = sizeof(*address);
    while (running)
    {
        pthread_t client;
        int client_fd = accept(server_fd, (struct sockaddr *)address, (socklen_t *)&addrlen);
        if (client_fd < 0)
        {
            perror("accept failed");
            close(server_fd);
            continue;
        }
        if (pthread_create(&client, NULL, handleClient, &client_fd) != 0)
        {
            perror("pthread_create failed");
            close(client_fd);
            continue;
        }

        pthread_detach(client);
    }
}

void *handleClient(void *arg)
{
    int client_fd = *(int *)arg;

    char buffer[BUFFER_SIZE] = {0};
    char *hello = "Hello from server";

    printf("Client %.3d entered the chat..\n", client_fd);

    while (1)
    {
        ssize_t bytes_read = read(client_fd, buffer, BUFFER_SIZE);
        if (bytes_read <= 0)
        {
            if (bytes_read == 0)
                printf("Client %.3d disconnected.\n", client_fd);
            else
                perror("read failed");
            break;
        }
        buffer[bytes_read] = '\0';
        printf("Client %.3d: %s\n", client_fd, buffer);
    }

    close(client_fd);
    printf("Client %.3d exited.\n", client_fd);

    pthread_exit(NULL);
}

void *handleTerminal(void *arg)
{
    char input[256];
    int client_fd;
    char message[256];

    while (running)
    {
        fgets(input, sizeof(input), stdin);
        input[strcspn(input, "\n")] = 0;

        if (sscanf(input, "%d : %[^\n]", &client_fd, message) == 2)
        {
            sendMessage(client_fd, message);
        }
        else
        {
            printf("Enter in format '<client_fd> : <message>\n");
        }
    }
    pthread_exit(NULL);
}

void sendMessage(int client_fd, const char message[])
{
    ssize_t bytes_sent = send(client_fd, message, strlen(message), 0);
    if (bytes_sent < 0)
    {
        perror("send failed");
    }
    else
    {
        printf("Message sent to client_fd %d: %s\n", client_fd, message);
    }
}
