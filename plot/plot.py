import matplotlib.pyplot as plt

# Data points for YCSB-A
x_a = [1, 2, 4, 8, 16, 32, 64]
y_a = [138.19, 163.06, 168.22, 168.01, 170.17, 171.75, 172.36]
x_a_d = [1, 2, 4, 8, 16, 32, 64]
y_a_d = [187.87, 62.94, 45.53, 43.08, 50.51, 53.36, 54.54]

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

# Data points for YCSB-B
x_b = [1, 2, 4, 8, 16, 32, 64]
y_b = [1151.76, 1447.76, 1629.73, 1683.84, 1735.37, 1729.37, 1776.40]
x_b_d = [1, 2, 4, 8, 16, 32, 64]
y_b_d = [1012.59, 891.78, 468.95, 493.71, 528.60, 557.69, 557.69]

# Creating the plot for YCSB-B
plt.figure(figsize=(10, 5))
plt.plot(x_b, y_b, marker='o', label='YCSB-B Normal')
plt.plot(x_b_d, y_b_d, marker='o', label='YCSB-B with Network Disconnection')
plt.title('Throughput vs. Number of Clients [YCSB-B]')
plt.xlabel('Number of Clients')
plt.ylabel('Throughput (ops/sec)')
plt.legend()
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/throughput_plot-B.png')
plt.close()

# Data points for YCSB-C
x_c = [1, 2, 4, 8, 16, 32, 64]
y_c = [2299.35, 4824.97, 7766.40, 13626.49, 19936.87, 22872.54, 34586.98]
x_c_d = [1, 2, 4, 8, 16, 32, 64]
y_c_d = [2490.84, 4968.88, 8424.45, 14066.70, 20734.94, 23600.02, 37699.64]

# Creating the plot for YCSB-C
plt.figure(figsize=(10, 5))
plt.plot(x_c, y_c, marker='o', label='YCSB-C Normal')
plt.plot(x_c_d, y_c_d, marker='o', label='YCSB-C with Network Disconnection')
plt.title('Throughput vs. Number of Clients [YCSB-C]')
plt.xlabel('Number of Clients')
plt.ylabel('Throughput (ops/sec)')
plt.legend()
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/throughput_plot-C.png')
plt.close()