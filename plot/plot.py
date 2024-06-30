import matplotlib.pyplot as plt

# Data points for YCSB-A
x_a = [1, 2, 4, 8, 16, 32, 64]
y_a= [259.17,407.36,461.94,457.64,461.21,480.75,483.78]


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

# Data points for YCSB-B
x_b = [1, 2, 4, 8, 16, 32, 64]

y_b = [763.34,1484.16,2767.86,4613.38,4932.94,4972.38,4917.45]




x_b_d = [1, 2, 4, 8, 16, 32, 64]
y_b_d= [749.08,1489.32,2446.26,1645.15,1725.82,1746.60,1856.93]



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
y_c = [944.33,1863.22,3653.79,7282.32,13235.30,19410.69,28014.23]

x_c_d = [1, 2, 4, 8, 16, 32, 64]
y_c_d = [968.66,1822.50,3636.08,6972.39,13143.13,19396.28,27571.92]


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