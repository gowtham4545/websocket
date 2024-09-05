#include <stdio.h>
#include "client.h"
#include <netinet/in.h>
#include <pthread.h>
#include <unistd.h>
#include <errno.h>

int main()
{
    pthread_t terminal;
    while (1)
    {
        int fd = createClient();
        if (pthread_create(&terminal, NULL, handleTerminal, &fd) != 0)
        {
            perror("Failed to create terminal thread");
            return 1;
        }
        if (!errno)
            communicateServer(fd);
        close(fd);
        printf("Reconnecting after 5 seconds...\n");
        sleep(5);
        errno = 0;
    }
    printf("Client terminated.\n");

    return 0;
}
