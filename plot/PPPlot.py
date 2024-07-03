import matplotlib as mpl
import matplotlib.pyplot as plt
import numpy as np

# フォントの設定
mpl.rcParams['font.family'] = 'Hiragino Sans'  
mpl.rcParams['font.size'] = 30  

# データの定義
throughput = {
    'normal': [(4020+3879+3919+3895+4127)/100, (6319+6417+6349+6156+6210)/100, (7217+6737+7266+7277+6599)/100],
    'partition': [(5981+6023+6177+6000+5892)/150, (9309+9983+9941+10074+9835)/150, (10449+10353+10879+11084+10797)/150]
}
latency = {
    'normal': [19.999/((4020+3879+3919+3895+4127)/5),39.995/((6319+6417+6349+6156+6210)/5),  60.065/((7217+6737+7266+7277+6599)/5)],
    'partition': [30.012/((5981+6023+6177+6000+5892)/5), 60.065/((9309+9983+9941+10074+9835)/5),89.997/((10449+10353+10879+11084+10797)/5)]
}

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

x_a = [1, 2, 3]
y_a = [(4020+3879+3919+3895+4127)/100, (6319+6417+6349+6156+6210)/100, (7217+6737+7266+7277+6599)]
y_a_lat = [19.999/((4020+3879+3919+3895+4127)/5),39.995/((6319+6417+6349+6156+6210)/5),  60.065/((7217+6737+7266+7277+6599)/5)]

x_a_crash = [1, 2, 3]
y_a_crash = [(5981+6023+6177+6000+5892)/150, (9309+9983+9941+10074+9835)/150, (10449+10353+10879+11084+10797)/150]
y_a_crash_lat = [30.012/((5981+6023+6177+6000+5892)/5), 60.065/((9309+9983+9941+10074+9835)/5),89.997/((10449+10353+10879+11084+10797)/5)]
