import matplotlib.pyplot as plt

# Data points for YCSB-A
x_a = [1, 2, 4, 8, 16, 32, 64]
y_a= []


x_a_d = [1, 2, 4, 8, 16, 32, 64]

y_a_d = [232.73,131.24,104.24,108.77,115.84,106.50,105.11]




# Creating the plot for YCSB-A
plt.figure(figsize=(10, 5))
plt.plot(x_a, y_a, marker='o', label='YCSB-A Normal')
plt.plot(x_a_d, y_a_d, marker='o', label='YCSB-A with Network Disconnection')
plt.title('Throughput vs. Number of Clients [YCSB-A]')
plt.xlabel('Number of Clients')
plt.ylabel('Throughput (ops/sec)')
plt.legend()
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/throughput_plot-A.png')
plt.close()