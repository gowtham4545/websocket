#include <stdio.h>
#include "client.h"
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <pthread.h>

#define BUFFER_SIZE 1024

int createClient()
{
    struct sockaddr_in serv_addr;
    int fd = createSocket();
    connectServer(fd, &serv_addr);
    return fd;
}

int createSocket()
{
    int fd = socket(AF_INET, SOCK_STREAM, 0);
    if (fd < 0)
    {
        printf("\nSocket creation error \n");
        exit(EXIT_FAILURE);
    }
    return fd;
}

void connectServer(int fd, struct sockaddr_in *serv_addr)
{
    serv_addr->sin_family = AF_INET;
    serv_addr->sin_port = htons(SERVER_PORT);

    if (inet_pton(AF_INET, "127.0.0.1", &serv_addr->sin_addr) <= 0)
    {
        printf("\nInvalid address/ Address not supported \n");
        exit(EXIT_FAILURE);
    }

    if (connect(fd, (struct sockaddr *)serv_addr, sizeof(*serv_addr)) < 0)
    {
        printf("\nConnection Failed \n");
        return;
    }
    printf("Connected to server...\n");
}

void communicateServer(int fd)
{
    char buffer[BUFFER_SIZE] = {0};

    while (1)
    {

        ssize_t bytes_read = read(fd, buffer, BUFFER_SIZE);
        if (bytes_read < 0)
        {
            perror("Read failed");
            break;
        }
        else if (bytes_read == 0)
        {
            printf("Server closed the connection.\n");
            break;
        }
        else
        {
            buffer[bytes_read] = '\0';
            printf("Message from server: %s\n", buffer);
        }
        memset(buffer, 0, BUFFER_SIZE);
    }
}

void *handleTerminal(void *arg)
{
    char input[256];
    int server_fd = *(int *)arg;

    while (1)
    {
        fgets(input, sizeof(input), stdin);
        input[strcspn(input, "\n")] = 0;

        if (strcmp(input, "exit") == 0)
            exit(0);

        sendMessage(server_fd, input);
        memset(input, '\0', 256);
    }
    pthread_exit(NULL);
}

void sendMessage(int fd, const char message[])
{
    ssize_t bytes_sent = send(fd, message, strlen(message), 0);
    if (bytes_sent < 0)
    {
        perror("send failed");
    }
}
