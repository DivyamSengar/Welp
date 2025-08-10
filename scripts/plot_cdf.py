import sys
import numpy as np
import matplotlib.pyplot as plt

# Function to extract latency in microseconds and convert to milliseconds
def extract_latency_in_ms(line):
    if "took" in line:
        try:
            # Extract the part after "took" and before "us"
            latency_us = int(line.split("took")[1].split("us")[0].strip())
            # Convert microseconds to milliseconds
            return latency_us / 1000.0
        except (IndexError, ValueError):
            return None
    return None

# List to store latency values in milliseconds
latencies_in_ms = []

with open(sys.argv[1], 'r') as f:
    log_lines = f.readlines()

# Extract latency values from the log lines
for line in log_lines:
    latency_ms = extract_latency_in_ms(line)
    if latency_ms is not None:
        latencies_in_ms.append(latency_ms)

# Sort latencies to calculate the CDF
latencies_in_ms_sorted = np.sort(latencies_in_ms)

# Calculate the CDF values
cdf = np.arange(1, len(latencies_in_ms_sorted) + 1) / len(latencies_in_ms_sorted)

# Plot the CDF
plt.step(latencies_in_ms_sorted, cdf, where='post', label="CDF")
plt.xlabel('Latency (ms)')
plt.ylabel('CDF')
plt.title('Cumulative Distribution Function of Latency (ms)')
plt.grid(True)
plt.xlim(left=0)
plt.savefig(sys.argv[1] + ".png")
