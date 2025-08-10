# Thread Affinity and Shared Memory Benchmark

This project demonstrates how to use threads with shared memory and measure the impact of thread affinity (i.e. whether threads run on the same CPU core or on different cores). The code also measures execution time for a simple workload and can optionally be run under Linux's `perf` tool to capture additional performance metrics.

## Files

- **thread_mini.c**  
  Contains the C source code for the benchmark. The program creates two threads that both work on a shared memory buffer, synchronizing updates with a mutex. Each thread runs the workload for a specified number of iterations and prints timing information.

- **benchmark.sh**  
  A bash script to compile the code and run the benchmark. It also supports running the binary under `perf` for additional metrics.

## Requirements

- **Linux:**  
  CPU affinity calls (using `CPU_SET`, `CPU_ZERO`, and `pthread_setaffinity_np`) are only supported on Linux. On non-Linux systems (e.g. macOS), the affinity calls will be skipped.
- **GCC:**  
  For compiling the C source code.
- **perf (optional):**  
  Linux performance monitoring tool. Install via your package manager (e.g., `sudo apt-get install linux-tools-common linux-tools-generic`).

## Command Line Options

When running the benchmark, you can pass the following options:

- `-a same` or `--affinity same`  
  Pin both threads to the same CPU (CPU 0).  
  Default is `diff` (threads on different cores: thread1 on CPU 0 and thread2 on CPU 1).

- `-i iterations` or `--iterations iterations`  
  Specify the number of iterations per thread.  
  Default is `1000`.

- `-p` or `--perf`  
  Print instructions for using `perf` and, if specified as the first argument in the bash script, run the benchmark under `perf`.

- `-h` or `--help`  
  Display usage information.

## How to Run

1. **Compile and Run Normally:**

   - Make the bash script executable:
     ```bash
     chmod +x benchmark.sh
     ```
   - Run the benchmark with default settings (threads on different cores, 1000 iterations):
     ```bash
     ./benchmark.sh
     ```

2. **Run with Threads on the Same Core:**

   ```bash
   ./benchmark.sh -a same
   ```

3. **Change the Number of Iterations:**

   For example, to run 5000 iterations per thread:
   ```bash
   ./benchmark.sh -i 5000
   ```

4. **Run Under Perf for Additional Metrics:**

   Prepend `--perf` before other options to run the benchmark under `perf`:
   ```bash
   ./benchmark.sh --perf -a diff -i 5000
   ```
   This will execute:
   ```bash
   perf stat -e cycles,cache-references,cache-misses ./my_benchmark -a diff -i 5000
   ```
   and print additional performance metrics.

## Expected Output

Each thread will output its total time and average time per iteration. For example:
```
Thread 1: Total time for 1000 iterations: 0.123456 seconds (avg: 0.000123 sec/iter)
Thread 2: Total time for 1000 iterations: 0.125678 seconds (avg: 0.000126 sec/iter)
Shared memory contents after business logic:
2 2 2 ... 2
```

When using `perf`, you will also see additional statistics such as CPU cycles, cache references, and cache misses.

## Notes

- On macOS or other non-Linux systems, thread affinity settings will not be applied, and a message will indicate that affinity is not supported.
- Adjust the workload and iterations to best suit your experimental needs.

---

