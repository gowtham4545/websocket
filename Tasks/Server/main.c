#include "stdio.h"
#include "server.h"
#include <netinet/in.h>
#include <unistd.h>
#include <pthread.h>
#include <signal.h>

int server_fd;
void handleSigint(int sig)
{
    running = 0;
    close(server_fd);
    exit(0);
}

int main()
{
    struct sockaddr_in address;
    pthread_t terminal;
    signal(SIGINT, handleSigint);

    server_fd = createServer();

    if (pthread_create(&terminal, NULL, handleTerminal, NULL) != 0)
    {
        perror("Failed to create terminal thread");
        return 1;
    }

    acceptConnection(server_fd, &address);

    close(server_fd);

    pthread_join(terminal, NULL);
    printf("Server terminated.\n");

    return 0;
}
