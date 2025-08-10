import matplotlib.pyplot as plt
import numpy as np

# Expected format: {througput: [a list of latencies]} (For throughput-latency graph, use average throughput)
def draw(latency_dict, title):
    fig, ax = plt.subplots(figsize=(10, 6))

    # Sorting the latency_dict keys
    rate = sorted(latency_dict.keys())
    y = []
    std = []

    for r in rate:
        latency = np.array(latency_dict[r])
        y.append(np.average(latency))
        std.append(np.std(latency))
        print(f"Throughput: {r}, Latencies: {latency}, Avg: {y[-1]}, Std: {std[-1]}")

    # Plotting error bars with better formatting
    ax.errorbar(rate, y, yerr=std, fmt='-o', color='b', ecolor='gray', elinewidth=1, capsize=5, label='99th Percentile Latency')

    # Add a horizontal line at y=500
    ax.axhline(y=500, color='r', linestyle='--', label='SLO Target (500 ms)')

    # Setting labels and title
    ax.set_xlabel("Throughput (Req/Sec)")
    ax.set_ylabel("Latency (ms)")
    ax.set_title(title)

    # Adding gridlines
    ax.grid(True, linestyle='--', alpha=0.7)

    # Adding a legend
    ax.legend(loc='upper left')

    # Customize the appearance
    ax.set_facecolor('#f0f0f0')  # Set background color
    ax.spines['top'].set_visible(False)
    ax.spines['right'].set_visible(False)
    ax.spines['bottom'].set_linewidth(0.5)
    ax.spines['left'].set_linewidth(0.5)
    ax.tick_params(axis='both', which='both', width=0.5)

    # Show the plot
    plt.show()

latency_dict = {
    100: [200, 210, 190, 205, 215],  # 100 Req/Sec, corresponding latencies
    200: [300, 290, 310, 295, 285],  # 200 Req/Sec, corresponding latencies
    300: [400, 410, 420, 395, 405],  # 300 Req/Sec, corresponding latencies
    400: [600, 590, 610, 605, 615]   # 400 Req/Sec, corresponding latencies
}

draw(latency_dict, 'Testa')
