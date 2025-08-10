#define _GNU_SOURCE
#include <stdio.h>
#include <sys/mman.h>
#include <sched.h>
#include <unistd.h>
#include <string.h>
#include <errno.h>
#include <pthread.h>
#include <stdlib.h>
#include <time.h>
#include <getopt.h>
#include <math.h>

#define BUF_SIZE 32

// Global mutex for synchronizing access to shared memory.
pthread_mutex_t mutex = PTHREAD_MUTEX_INITIALIZER;

// Business logic: increments each byte of the shared buffer.
void do_business_logic(int buf[]) {
    pthread_mutex_lock(&mutex);
    if (buf[0] == NULL) {
        fprintf(stderr, "ERROR: business_logic: buf is empty!\n");
        pthread_mutex_unlock(&mutex);
        return;
    }
    for (int i = 0; i < BUF_SIZE; i++) {
        buf[i] += 15;
        buf[i] /= 71;
        buf[i] = buf[i]^2;
        // buf[i] = sin(buf[i]);
        // float val = (float)buf[i];
        // buf[i] = (int)(val + sin(i * 0.1) * 10); // Using sin function
    }
    pthread_mutex_unlock(&mutex);
}

// Structure to pass arguments to each thread.
typedef struct {
    int *buf;
    int iterations;
    int thread_id;
} thread_data_t;

// Each thread runs this function: loops over iterations and times the work.
void *thread_work(void *arg) {
    thread_data_t *data = (thread_data_t *) arg;
    struct timespec start, end;
    clock_gettime(CLOCK_MONOTONIC, &start);
    for (int i = 0; i < data->iterations; i++) {
        do_business_logic(data->buf);
    }
    clock_gettime(CLOCK_MONOTONIC, &end);
    double elapsed = (end.tv_sec - start.tv_sec) + (end.tv_nsec - start.tv_nsec) / 1e9;
    printf("Thread %d: Total time for %d iterations: %f seconds (avg: %f sec/iter)\n",
           data->thread_id, data->iterations, elapsed, elapsed / data->iterations);
    return NULL;
}

void print_usage(char *progname) {
    printf("Usage: %s [-a same|diff] [-i iterations] [--perf]\n", progname);
    printf("  -a, --affinity   Set thread affinity: 'same' to run both threads on CPU 0, 'diff' to run on CPU 0 and CPU 1 (default: diff)\n");
    printf("  -i, --iterations Number of iterations per thread (default: 1000)\n");
    printf("  -p, --perf       Print instructions for using perf for additional metrics\n");
    printf("  -h, --help       Display this help message\n");
}

int main(int argc, char *argv[]) {
    int opt;
    int iterations = 1000;
    int use_perf = 0;
    char affinity_option[16] = "diff";  // Default: threads on different cores.

    static struct option long_options[] = {
        {"affinity", required_argument, 0, 'a'},
        {"iterations", required_argument, 0, 'i'},
        {"perf", no_argument, 0, 'p'},
        {"help", no_argument, 0, 'h'},
        {0, 0, 0, 0}
    };

    // Parse command-line options.
    while ((opt = getopt_long(argc, argv, "a:i:ph", long_options, NULL)) != -1) {
        switch (opt) {
            case 'a':
                strncpy(affinity_option, optarg, sizeof(affinity_option) - 1);
                affinity_option[sizeof(affinity_option) - 1] = '\0';
                break;
            case 'i':
                iterations = atoi(optarg);
                break;
            case 'p':
                use_perf = 1;
                break;
            case 'h':
            default:
                print_usage(argv[0]);
                exit(EXIT_SUCCESS);
        }
    }

    if (use_perf) {
        printf("To measure additional metrics with perf, run the binary using:\n");
        printf("   perf stat -e cycles,cache-references,cache-misses ./%s", argv[0]);
        for (int i = 1; i < argc; i++) {
            printf(" %s", argv[i]);
        }
        printf("\n");
    }

    // Allocate shared memory.
    int *detail_buf = mmap(NULL, sizeof(int) * BUF_SIZE, PROT_WRITE | PROT_READ,
                            MAP_SHARED | MAP_ANONYMOUS, -1, 0);
    if (detail_buf == MAP_FAILED) {
        fprintf(stderr, "mmap failed, errno: %d\n", errno);
        return 1;
    }


    // Initialize the shared buffer with nonzero values.
    memset(detail_buf, 1, BUF_SIZE);

    pthread_t thread1, thread2;
    thread_data_t data1 = { detail_buf, iterations, 1 };
    thread_data_t data2 = { detail_buf, iterations, 2 };

    // Create threads.
    if (pthread_create(&thread1, NULL, thread_work, &data1) != 0) {
        perror("pthread_create thread1");
        exit(EXIT_FAILURE);
    }
    if (pthread_create(&thread2, NULL, thread_work, &data2) != 0) {
        perror("pthread_create thread2");
        exit(EXIT_FAILURE);
    }

    // Set thread affinity.
    cpu_set_t cpuset;
    // For thread1, always pin to CPU 0.
    CPU_ZERO(&cpuset);
    CPU_SET(0, &cpuset);
    if (pthread_setaffinity_np(thread1, sizeof(cpu_set_t), &cpuset) != 0) {
        perror("pthread_setaffinity_np thread1");
    }

    // For thread2, use the command-line option.
    CPU_ZERO(&cpuset);
    if (strcmp(affinity_option, "same") == 0) {
        CPU_SET(0, &cpuset);  // Both on CPU 0.
    } else {
        CPU_SET(1, &cpuset);  // Default: thread2 on CPU 1.
    }
    if (pthread_setaffinity_np(thread2, sizeof(cpu_set_t), &cpuset) != 0) {
        perror("pthread_setaffinity_np thread2");
    }

    // Wait for threads to finish.
    pthread_join(thread1, NULL);
    pthread_join(thread2, NULL);

    // Display final shared memory state.
    printf("Shared memory contents after business logic:\n");
    for (int i = 0; i < BUF_SIZE; i++) {
        printf("%d ", detail_buf[i]);
    }
    printf("\n");

    return 0;
}
