import matplotlib.pyplot as plt


fontSize = 14

# Data points for YCSB-A
x_a = [1, 2, 4]
y_a= [(5565+5322+5208)/90,(8411+8182+8288)/90 ,(9780+9378+9268)/90]

y_a_lat = [29.99356, 59.992, 120.0688 ]


x_a_d = [1, 2, 4 ]

y_a_d = [(5438+5353+5265)/90,(3598+3202+3466)/90,(2883+3033+2844)/90]

y_a_d_lat = [29.992,60.0031,120.0489]






# Creating the plot for YCSB-A
plt.figure(figsize=(10, 5))
plt.plot(x_a, y_a, marker='o', label='YCSB-A Normal')
plt.plot(x_a_d, y_a_lat, marker='o', label='YCSB-A with Network Disconnection')
plt.xlabel('Number of Clients')
plt.ylabel('Throughput (ops/sec)')
plt.legend(fontsize=fontSize)
plt.rcParams["font.size"] = 14
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/throughput_plot-A.png')
plt.close()

# Creating the plot for YCSB-A
plt.figure(figsize=(10, 5))
plt.plot(x_a, y_a_lat, marker='o', label='YCSB-A Normal')
plt.plot(x_a_d, y_a_d_lat, marker='o', label='YCSB-A with Network Disconnection')
plt.xlabel('Number of Clients')
plt.ylabel('Latency for each operation (ms)')
plt.legend(fontsize=fontSize)
plt.rcParams["font.size"] = 14
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/Latency-A.png')
plt.close()


x_a = [1, 2, 4]
y_a= [(5565+5322+5208)/90,(8411+8182+8288)/90 ,(9780+9378+9268)/90]


x_a = [1, 2, 4]
y_a_crash = [(5436+5448+5008)/90,(7638+7979+8446)/90, (8673+9129+8514)/90]
y_a_crash_lat = [29.9972,60.0035,120.033]

plt.figure(figsize=(10, 5))
plt.plot(x_a, y_a, marker='o', label='YCSB-A Normal')
plt.plot(x_a_d, y_a_crash, marker='o', label='YCSB-A with Replica Crash')
plt.title('Throughput vs. Number of Clients [YCSB-A]')
plt.xlabel('Number of Clients')
plt.ylabel('Throughput (ops/sec)')
plt.legend(fontsize=fontSize)
plt.rcParams["font.size"] = 14
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/throughput_with_crash_plot-A.png')
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


