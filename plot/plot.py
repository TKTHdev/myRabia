import matplotlib.pyplot as plt


fontSize = 14

x_a = [1, 2, 3]
y_a = [178.83, 276.46, 315.84]
y_a_lat = [0.005603, 0.007221, 0.012308]
x_a_d = [1, 2, 3]
y_a_d = [178.40, 113.96, 109.42]
y_a_d_lat = [0.005603, 0.017519, 0.036601]



# Creating the plot for YCSB-A
plt.figure(figsize=(10, 5))
plt.plot(x_a, y_a, marker='o', label='YCSB-A Normal')
plt.plot(x_a_d, y_a_d, marker='o', label='YCSB-A with Network Disconnection')
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



x_a = [1, 2, 3]
y_a = [(4020+3879+3919+3895+4127)/100, (6319+6417+6349+6156+6210)/100, (7217+6737+7266+7277+6599)]
y_a_lat = [19.999/((4020+3879+3919+3895+4127)/5),39.995/((6319+6417+6349+6156+6210)/5),  60.065/((7217+6737+7266+7277+6599)/5)]

x_a_crash = [1, 2, 3]
y_a_crash = [(5981+6023+6177+6000+5892)/150, (9309+9983+9941+10074+9835)/150, (10449+10353+10879+11084+10797)/150]
y_a_crash_lat = [30.012/((5981+6023+6177+6000+5892)/5), 60.065/((9309+9983+9941+10074+9835)/5),89.997/((10449+10353+10879+11084+10797)/5)]

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

plt.figure(figsize=(10, 5))
plt.plot(x_a, y_a_lat, marker='o', label='YCSB-A Normal')
plt.plot(x_a_d, y_a_crash_lat, marker='o', label='YCSB-A with Replica Crash')
plt.title('Latency vs. Number of Clients [YCSB-A]')
plt.xlabel('Number of Clients')
plt.ylabel('Latency for each operation (ms)')
plt.legend(fontsize=fontSize)
plt.rcParams["font.size"] = 14
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/latency-A.png')
plt.close()


# Data points for YCSB-C
x_c = [1, 2, 3]
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


x_b = [1, 2, 3]
y_b = [(19849+18120+18312+19884+20320)/150, (39155+37878+40574+40981+41650)/150, (51998+52372+56123+54938+59552)/150]
y_b_lat = [30/((19849+18120+18312+19884+20320)/5), 60/((39155+37878+40574+40981+41650)/5), 90/((51998+52372+56123+54938+59552)/5)]


x_c = [1, 2, 3]
y_c =[(25600+25602+27349+26022+24428)/150 , (52150+53915+52436+50504+49684)/150, (75316+74925+77946+79475+78959)/150]
y_c_lat = [30/((25600+25602+27349+26022+24428)/5), 60/((52150+53915+52436+50504+49684)/5), 90/((75316+74925+77946+79475+78959)/5)]


x_5 = [1, 2, 3]
y_5 = [(4306+4362+4449)/90, (6315+5695+6551)/90, (6739+6941+7016+6659+7053)/120]
y_lat = [30/((4306+4362+4449)/3), 60/((6315+5695+6551)/3), 90/((6739+6941+7016+6659+7053)/5)]

x_7= [1, 2, 3]
y_7 = [(3089+3492+3251+3242)/120, (4453+4615+4276+4593)/120, (4598+5005+4979+4783+4855)/120]
y_7_lat = [30/((3089+3492+3251+3242)/4), 60/((4453+4615+4276+4593)/4), 90/((4598+5005+4979+4783+4855)/5)]


x_geo = [1, 2, 3]
y_geo = [(23+23+23)/90, (39+43+51+45+44+43)/150,(55+46+49+64)/120]
y_geo_lat = [30/((23+23+23)/3), 60/((39+43+51+45+44+43)/6), 90/((55+46+49+64)/4)]