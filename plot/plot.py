import matplotlib.pyplot as plt


fontSize = 14

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
plt.legend(fontsize=fontSize)
plt.rcParams["font.size"] = 14
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
plt.legend(fontsize=fontSize)
plt.rcParams["font.size"] = fontSize
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
plt.legend(fontsize=fontSize)
plt.rcParams["font.size"] = fontSize
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/throughput_plot-C.png')
plt.close()



x_a = [1, 2, 4, 8, 16, 32, 64]
y_a= [259.17,407.36,461.94,457.64,461.21,480.75,483.78]
x_a_geo = [1, 2, 4, 8, 16, 32, 64]
y_a_geo = [1.04,1.40,1.72,1.88,2.68,3.42,5.65]
x_b = [1, 2, 4, 8, 16, 32, 64]
y_b = [763.34,1484.16,2767.86,4613.38,4932.94,4972.38,4917.45]
x_b_geo = [1, 2, 4, 8, 16, 32, 64]
y_b_geo = [5.70,9.52,14.17,18.17,23.29,33.78,51.71]
x_c = [1, 2, 4, 8, 16, 32, 64]
y_c = [944.33,1863.22,3653.79,7282.32,13235.30,19410.69,28014.23]

x_c_geo = [1, 2, 4, 8, 16, 32, 64]
y_c_geo = [9.42,16.22,29.32,57.26,110.58,219.58,434.26]



plt.figure(figsize=(10, 5))
#plt.plot(x_a, y_a, marker='o', label='YCSB-A Normal')
plt.plot(x_a_geo, y_a_geo, marker='o', label='YCSB-A when Geo-Distributed')
plt.plot(x_b_geo, y_b_geo, marker='o', label='YCSB-B when Geo-Distributed')
plt.plot(x_c_geo, y_c_geo, marker='o', label='YCSB-C when Geo-Distributed')
plt.title('Throughput vs. Number of Clients [YCSB-A,B,C]')
plt.xlabel('Number of Clients')
plt.ylabel('Throughput (ops/sec)')
plt.legend(fontsize=fontSize)
plt.rcParams["font.size"] = fontSize
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/throughput_plot-_geo.png')
plt.close()


x_b = [1, 2, 4, 8, 16, 32, 64]
y_b = [763.34,1484.16,2767.86,4613.38,4932.94,4972.38,4917.45]
x_b_geo = [1, 2, 4, 8, 16, 32, 64]
y_b_geo = [5.70,9.52,14.17,18.17,23.29,33.78,51.71]

plt.figure(figsize=(10, 5))
#plt.plot(x_b, y_b, marker='o', label='YCSB-B Normal')
plt.plot(x_b_geo, y_b_geo, marker='o', label='YCSB-B when Geo-Distributed')
plt.title('Throughput vs. Number of Clients [YCSB-B]')
plt.xlabel('Number of Clients')
plt.ylabel('Throughput (ops/sec)')
plt.legend(fontsize=fontSize)
plt.rcParams["font.size"] = fontSize
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/throughput_plot-B_geo.png')
plt.close()


x_c = [1, 2, 4, 8, 16, 32, 64]
y_c = [944.33,1863.22,3653.79,7282.32,13235.30,19410.69,28014.23]

x_c_geo = [1, 2, 4, 8, 16, 32, 64]
y_c_geo = [9.42,16.22,29.32,57.26,110.58,219.58,434.26]

plt.figure(figsize=(10, 5))
#plt.plot(x_c, y_c, marker='o', label='YCSB-C Normal')
plt.plot(x_c_geo, y_c_geo, marker='o', label='YCSB-C when Geo-Distributed')
plt.title('Throughput vs. Number of Clients [YCSB-C]')
plt.xlabel('Number of Clients')
plt.ylabel('Throughput (ops/sec)')
plt.legend(fontsize=fontSize)
plt.rcParams["font.size"] = fontSize
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/throughput_plot-C_geo.png')
plt.close()




x_o = [1, 2, 3, 4,5,6,7,8]
y_o = [61695,112742,124189,145647,117635,126337, 120289, 109782]


plt.figure(figsize=(10, 5))
#plt.plot(x_c, y_c, marker='o', label='YCSB-C Normal')
plt.plot(x_o, y_o, marker='o', label='YCSB-C when Geo-Distributed')
plt.title('Throughput vs. Number of Clients [YCSB-C]')
plt.xlabel('Number of Clients')
plt.ylabel('Throughput (ops/sec)')
plt.legend(fontsize=fontSize)
plt.rcParams["font.size"] = fontSize
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/throughput_plo_geo.png')
plt.close()



