import matplotlib.pyplot as plt

# Data collected from the Go script (example data)
data = {
    'MessageCount': [10, 50, 100, 200, 500,2000,5000,10000],
    'AverageLatency_ms': [247.32, 9.08, 54.99, 63.08, 63.91,248.60,240.29,422.48],
    'Throughput_msg_sec': [40.37,2294.58,1700.85,2813.87,4393.50,6269.41,14453.17,17142.55]
}

# Plotting Average Latency
plt.figure(figsize=(10, 5))
plt.plot(data['MessageCount'], data['AverageLatency_ms'], marker='o')
plt.title('Average Latency vs. Message Count')
plt.xlabel('Message Count')
plt.ylabel('Average Latency (ms)')
plt.grid(True)
plt.show()

# Plotting Throughput
plt.figure(figsize=(10, 5))
plt.plot(data['MessageCount'], data['Throughput_msg_sec'], marker='o', color='r')
plt.title('Throughput vs. Message Count')
plt.xlabel('Message Count')
plt.ylabel('Throughput (msg/sec)')
plt.grid(True)
plt.show()
