#define _GNU_SOURCE
#include <stdio.h>
#include <sys/mman.h>
#include <sched.h>
#include <unistd.h>
#include <string.h>
#include <errno.h>
#include <pthread.h>
#define buf_size 32

enum workloads {
    detail,
    reservation,
    review
};

int main (int argc, char* argv[]) {
    // Create buffer in shared memory
    void* detail_buf = mmap(NULL, buf_size, PROT_WRITE | PROT_READ, MAP_SHARED | MAP_ANONYMOUS, -1, 0);
    // Check for errors
    if (detail_buf == MAP_FAILED) {
        printf("mmap failed, errno: %d\n", errno);
        return 1;
    }
    // One initialize the shared memory
    memset(detail_buf, 1, buf_size);
    // Create a child process
    int child_pid = fork();
    // Schedule processes on the same CPU
    schedule_shared_procs(child_pid);
    // Do some work on the shared memory
    do_business_logic(detail_buf);
}

void schedule_shared_procs(int child_pid) {
    cpu_set_t mask;
    CPU_ZERO(&mask);
    CPU_SET(0, &mask);
    // Set the parent and child to run on CPU 0
    sched_setaffinity(getpid(), sizeof(cpu_set_t), &mask);
    sched_setaffinity(child_pid, sizeof(cpu_set_t), &mask);
}

void do_business_logic (char buf[]) {
    // Check buffer nonempty
    if ((buf[0]) == 0) {
        printf("ERROR: business_logic: buf is empty!\n");
        return;
    }
    // Increment entries
    for (int i = 0; i < buf_size; i++) {
        buf[i] += 1;
    }
}

void *thread_work(void *arg) {
    // threads run the business logic code here
    do_business_logic((char *)arg);
    return NULL;
}