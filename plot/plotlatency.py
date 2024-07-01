import matplotlib.pyplot as plt

# Data points for YCSB-A
x_a = [1, 2, 4, 8, 16, 32, 64]
Aread_time = [1.167986, 1.110566, 1.177052, 1.160599, 1.173142, 1.166381, 1.248513]
Awrite_time = [6.062137, 7.282872, 13.123879, 27.224526, 55.740565, 112.683822, 226.901277]
Atotal_time = [3.625986, 4.213716, 7.121323, 14.154026, 28.570923, 56.93856, 115.006961]






# Creating the plot for YCSB-A
plt.figure(figsize=(10, 5))
plt.plot(x_a, Awrite_time, marker='o', label='Write latency')
plt.title('Latency vs. Number of Clients [YCSB-A]')
plt.xlabel('Number of Clients')
plt.ylabel('Latency for Write operation(Millsec/op)')
plt.legend()
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/latencyA/Wlatency_plot-A.png')
plt.close()

plt.figure(figsize=(10, 5))
plt.plot(x_a, Aread_time, marker='o', label='Read latency')
plt.title('Latency vs. Number of Clients [YCSB-A]')
plt.xlabel('Number of Clients')
plt.ylabel('Latency for Read operation(Millsec/op)')
plt.legend()
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/latencyA/Rlatency_plot-A.png')
plt.close()

x_b = [1, 2, 4, 8, 16, 32, 64]
Bread_time = [
    1.153677,
    1.127185,
    1.163875,
    1.187547,
    1.208233,
    1.210139,
    1.220473
]

Bwrite_time = [
    6.458360,
    6.711459,
    7.661122,
    12.746564,
    40.487863,
    101.125989,
    217.936855
]

Btotal_time = [
    1.425998,
    1.409084,
    1.484194,
    1.760910,
    3.168614,
    6.220303,
    12.082078
]


# Creating the plot for YCSB-A
plt.figure(figsize=(10, 5))
plt.plot(x_b, Bread_time, marker='o', label='Read latency')
plt.title('Read Latency vs. Number of Clients [YCSB-B]')
plt.xlabel('Number of Clients')
plt.ylabel('Latency for Read operation(Millsec/op)')
plt.legend()
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/latencyB/Rlatency_plot-B.png')
plt.close()

plt.figure(figsize=(10, 5))
plt.plot(x_b, Bwrite_time, marker='o', label='Write latency')
plt.title('Write Latency vs. Number of Clients [YCSB-B]')
plt.xlabel('Number of Clients')
plt.ylabel('Latency for Write operation(Millsec/op)')
plt.legend()
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/latencyB/Wlatency_plot-B.png')
plt.close()


x_c = [1, 2, 4, 8, 16, 32, 64]

Cread_time = [
    1.150339,
    1.096976,
    1.114770,
    1.135267,
    1.251066,
    1.659278,
    2.410071
]

plt.figure(figsize=(10, 5))
plt.plot(x_c, Cread_time, marker='o', label='Read latency')
plt.title('Read Latency vs. Number of Clients [YCSB-C]')
plt.xlabel('Number of Clients')
plt.ylabel('Latency for Read operation(Millsec/op)')
plt.legend()
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/latencyC/Rlatency_plot-C.png')
plt.close()


