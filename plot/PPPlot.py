import matplotlib as mpl
import matplotlib.pyplot as plt
import numpy as np

# フォントの設定
mpl.rcParams['font.family'] = 'Hiragino Sans'  
mpl.rcParams['font.size'] = 30  

# データの定義
throughput = {
    'normal': [(4020+3879+3919+3895+4127)/100, (6319+6417+6349+6156+6210)/100, (7217+6737+7266+7277+6599)/100],
    'replica crash': [(5981+6023+6177+6000+5892)/150, (9309+9983+9941+10074+9835)/150, (10449+10353+10879+11084+10797)/150]
}
latency = {
    'normal': [19.999/((4020+3879+3919+3895+4127)/5),39.995/((6319+6417+6349+6156+6210)/5),  60.065/((7217+6737+7266+7277+6599)/5)],
    'replica crash': [30.012/((5981+6023+6177+6000+5892)/5), 60.065/((9309+9983+9941+10074+9835)/5),89.997/((10449+10353+10879+11084+10797)/5)]
}

plt.figure(figsize=(12, 8))  

plt.plot(throughput['normal'], latency['normal'], 'b-', label='通常', marker='o', markersize=10)
plt.plot(throughput['replica crash'], latency['replica crash'], 'g-', label='レプリカ故障あり', marker='s', markersize=10)

# グラフの設定
plt.xlabel('Throughput (reqs/sec)', fontsize=24)
plt.ylabel('Median Latency (ms)', fontsize=24)
plt.legend(fontsize=30)
plt.grid(True)

# 軸のフォントサイズを設定
plt.tick_params(axis='both', which='major', labelsize=12)

plt.tight_layout()

plt.savefig('replica_crash.png', dpi=300, bbox_inches='tight')

plt.show()

x_a = [1, 2, 3]
y_a = [(4020+3879+3919+3895+4127)/100, (6319+6417+6349+6156+6210)/100, (7217+6737+7266+7277+6599)]
y_a_lat = [19.999/((4020+3879+3919+3895+4127)/5),39.995/((6319+6417+6349+6156+6210)/5),  60.065/((7217+6737+7266+7277+6599)/5)]

x_a_crash = [1, 2, 3]
y_a_crash = [(5981+6023+6177+6000+5892)/150, (9309+9983+9941+10074+9835)/150, (10449+10353+10879+11084+10797)/150]
y_a_crash_lat = [30.012/((5981+6023+6177+6000+5892)/5), 60.065/((9309+9983+9941+10074+9835)/5),89.997/((10449+10353+10879+11084+10797)/5)]





# データの定義
throughput = {
    'normal': [178.83, 276.46, 315.84],
    'partition': [178.40, 113.96, 109.42]
}
latency = {
    'normal': [0.005603, 0.007221, 0.012308],
    'partition': [0.005603, 0.017519, 0.036601]
}

x_a = [1, 2, 3]
y_a = [178.83, 276.46, 315.84]
y_a_lat = [0.005603, 0.007221, 0.012308]
x_a_d = [1, 2, 3]
y_a_d = [178.40, 113.96, 109.42]
y_a_d_lat = [0.005603, 0.017519, 0.036601]


plt.figure(figsize=(12, 8))  

plt.plot(throughput['normal'], latency['normal'], 'b-', label='通常', marker='o', markersize=10)
plt.plot(throughput['partition'], latency['partition'], 'g-', label='ネットワーク分断あり', marker='s', markersize=10)

# グラフの設定
plt.xlabel('Throughput (reqs/sec)', fontsize=24)
plt.ylabel('Median Latency (ms)', fontsize=24)
plt.legend(fontsize=30)
plt.grid(True)

# 軸のフォントサイズを設定
plt.tick_params(axis='both', which='major', labelsize=12)

plt.tight_layout()

plt.savefig('network_partition.png', dpi=300, bbox_inches='tight')

plt.show()

# データの定義
throughput = {
    'YCSB-A': [178.83, 276.46, 315.84],
    'YCSB-B': [(19849+18120+18312+19884+20320)/150, (39155+37878+40574+40981+41650)/150, (51998+52372+56123+54938+59552)/150],
    'YCSB-C': [(25600+25602+27349+26022+24428)/150 , (52150+53915+52436+50504+49684)/150, (75316+74925+77946+79475+78959)/150]
    
}
latency = {
    'YCSB-A': [0.005603, 0.007221, 0.012308],
    'YCSB-B': [30/((19849+18120+18312+19884+20320)/5), 60/((39155+37878+40574+40981+41650)/5), 90/((51998+52372+56123+54938+59552)/5)],
    'YCSB-C': [30/((25600+25602+27349+26022+24428)/5), 60/((52150+53915+52436+50504+49684)/5), 90/((75316+74925+77946+79475+78959)/5)]
}



plt.figure(figsize=(12, 8))  

plt.plot(throughput['YCSB-A'], latency['YCSB-A'], 'b-', label='YCSB-A', marker='o', markersize=10)
plt.plot(throughput['YCSB-B'], latency['YCSB-B'], 'g-', label='YCSB-B', marker='s', markersize=10)
plt.plot(throughput['YCSB-C'], latency['YCSB-C'], 'c-', label='YCSB-C', marker='^', markersize=10)

# グラフの設定
plt.xlabel('Throughput (reqs/sec)', fontsize=24)
plt.ylabel('Median Latency (ms)', fontsize=24)
plt.legend(fontsize=30)
plt.grid(True)

# 軸のフォントサイズを設定
plt.tick_params(axis='both', which='major', labelsize=12)

plt.tight_layout()

plt.savefig('different benchmarks', dpi=300, bbox_inches='tight')

plt.show()

# データの定義
throughput = {
    'three': [178.83, 276.46, 315.84],
    'five': [(4306+4362+4449)/90, (6315+5695+6551)/90, (6739+6941+7016+6659+7053)/120],
    'seven': [(3089+3492+3251+3242)/120, (4453+4615+4276+4593)/120, (4598+5005+4979+4783+4855)/120]
    
}
latency = {
    'three': [0.005603, 0.007221, 0.012308],
    'five': [30/((3089+3492+3251+3242)/4), 60/((4453+4615+4276+4593)/4), 90/((4598+5005+4979+4783+4855)/5)],
    'seven': [30/((3089+3492+3251+3242)/4), 60/((4453+4615+4276+4593)/4), 90/((4598+5005+4979+4783+4855)/5)]
}



plt.figure(figsize=(12, 8))  

plt.plot(throughput['three'], latency['three'], 'b-', label='three', marker='o', markersize=10)
plt.plot(throughput['five'], latency['five'], 'g-', label='five', marker='s', markersize=10)
plt.plot(throughput['seven'], latency['seven'], 'r-', label='seven', marker='^', markersize=10)

# グラフの設定
plt.xlabel('Throughput (reqs/sec)', fontsize=24)
plt.ylabel('Median Latency (ms)', fontsize=24)
plt.legend(fontsize=30)
plt.grid(True)

# 軸のフォントサイズを設定
plt.tick_params(axis='both', which='major', labelsize=12)

plt.tight_layout()

plt.savefig('increasing_replicas.png', dpi=300, bbox_inches='tight')

plt.show()

x_5 = [1, 2, 3]
y_5 = [(4306+4362+4449)/90, (6315+5695+6551)/90, (6739+6941+7016+6659+7053)/120]
y_lat = [30/((4306+4362+4449)/3), 60/((6315+5695+6551)/3), 90/((6739+6941+7016+6659+7053)/5)]

x_7= [1, 2, 3]
y_7 = [(3089+3492+3251+3242)/120, (4453+4615+4276+4593)/120, (4598+5005+4979+4783+4855)/120]
y_7_lat = [30/((3089+3492+3251+3242)/4), 60/((4453+4615+4276+4593)/4), 90/((4598+5005+4979+4783+4855)/5)]


# データの定義
throughput = {
    'normal': [178.83, 276.46, 315.84],
    'geo': [(23+23+23)/90, (39+43+51+45+44+43)/150,(55+46+49+64)/120]
}
latency = {
    'normal': [0.005603, 0.007221, 0.012308],
    'geo': [30/((23+23+23)/3), 60/((39+43+51+45+44+43)/6), 90/((55+46+49+64)/4)]
}

plt.figure(figsize=(12, 8))  

plt.plot(throughput['normal'], latency['normal'], 'b-', label='通常', marker='o', markersize=10)
plt.plot(throughput['geo'], latency['geo'], 'g-', label='地理分散 ', marker='s', markersize=10)

# グラフの設定
plt.xlabel('Throughput (reqs/sec)', fontsize=24)
plt.ylabel('Median Latency (ms)', fontsize=24)
plt.legend(fontsize=30)
plt.grid(True)

# 軸のフォントサイズを設定
plt.tick_params(axis='both', which='major', labelsize=12)

plt.tight_layout()

plt.savefig('network_partition.png', dpi=300, bbox_inches='tight')

plt.show()

