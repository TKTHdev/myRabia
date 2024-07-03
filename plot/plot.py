import matplotlib.pyplot as plt


fontSize = 14

# Data points for YCSB-A
x_a = [1, 2, 4]
y_a= [259.17,407.36,461.94]


x_a_d = [1, 2, 4 ]

y_a_d = [232.73,131.24,104.24]




# Creating the plot for YCSB-A
plt.figure(figsize=(10, 5))
plt.plot(x_a, y_a, marker='o', label='YCSB-A Normal')
plt.plot(x_a_d, y_a_d, marker='o', label='YCSB-A with Network Disconnection')
plt.title('Throughput vs. Number of Clients [YCSB-A]')
plt.xlabel('Number of Clients')
plt.ylabel('Throughput (ops/sec)')
plt.legend(fontsize=fontSize)
plt.rcParams["font.size"] = 14
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/throughput_plot-A.png')
plt.close()

# Data points for YCSB-B
x_b = [1, 2, 4 ]

y_b = [763.34,1484.16,2767.86]






# Creating the plot for YCSB-B
plt.figure(figsize=(10, 5))
plt.plot(x_b, y_b, marker='o', label='YCSB-B Normal')
plt.title('Throughput vs. Number of Clients [YCSB-B]')
plt.xlabel('Number of Clients')
plt.ylabel('Throughput (ops/sec)')
plt.legend(fontsize=fontSize)
plt.rcParams["font.size"] = fontSize
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/throughput_plot-B.png')
plt.close()

# Data points for YCSB-C
x_c = [1, 2, 4 ]
y_c = [944.33,1863.22,3653.79]


# Creating the plot for YCSB-C
plt.figure(figsize=(10, 5))
plt.plot(x_c, y_c, marker='o', label='YCSB-C Normal')
plt.title('Throughput vs. Number of Clients [YCSB-C]')
plt.xlabel('Number of Clients')
plt.ylabel('Throughput (ops/sec)')
plt.legend(fontsize=fontSize)
plt.rcParams["font.size"] = fontSize
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/throughput_plot-C.png')
plt.close()



x_a = [1, 2, 4]
y_a= [259.17,407.36,461.94]
x_a_geo = [1, 2, 4]
y_a_geo = [1.04,1.40,1.72]



plt.figure(figsize=(10, 5))
#plt.plot(x_a, y_a, marker='o', label='YCSB-A Normal')
plt.plot(x_a, y_a, marker='o', label='YCSB-A Normal')
plt.plot(x_a_geo, y_a_geo, marker='o', label='YCSB-A when Geo-Distributed')
plt.title('Throughput vs. Number of Clients [YCSB-A,B,C]')
plt.xlabel('Number of Clients')
plt.ylabel('Throughput (ops/sec)')
plt.legend(fontsize=fontSize)
plt.rcParams["font.size"] = fontSize
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/throughput_plot-_geo.png')
plt.close()


