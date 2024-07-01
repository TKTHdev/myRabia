import matplotlib.pyplot as plt
import japanize_matplotlib

# Data points for YCSB-A
x_a = [1, 2, 4, 8, 16, 32, 64]
Aread_time = [1.288861, 1.351835, 1.400460, 1.396314, 1.368487, 1.354501, 1.462627]
Awrite_time = [9.006404, 11.966047, 23.411631, 48.891073, 96.797480, 188.424857, 402.058237]
Atotal_time = [10.295265, 13.317882, 24.811091, 50.287387, 98.165967, 189.779358, 403.510864]

# Creating the plot for YCSB-A
plt.figure(figsize=(10, 5))
plt.plot(x_a, Atotal_time, marker='o', label='latency')
for i, txt in enumerate(x_a):
    plt.annotate(txt, (x_a[i], Atotal_time[i]), xytext=(0, 5), textcoords='offset points', ha='center')
plt.xlabel('クライアントの数')
plt.ylabel('YCSB-Aにおける\n一命令あたりのレイテンシー(Millsec/命令)')
plt.legend()
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/latencyA/latency_plot-A.png')
plt.close()

# Data points for YCSB-B
x_b = [1, 2, 4, 8, 16, 32, 64]
Bread_time = [1.234441, 1.255589, 1.265219, 1.331707, 1.336240, 1.405194, 1.441195]
Bwrite_time = [8.192606, 9.434817, 12.405131, 24.326760, 48.578278, 159.547959, 355.424451]
Btotal_time = [9.427047, 10.690406, 13.670350, 25.658467, 49.914518, 160.953153, 356.865646]

# Creating the plot for YCSB-B
plt.figure(figsize=(10, 5))
plt.plot(x_b, Btotal_time, marker='o', label='latency')
for i, txt in enumerate(x_b):
    plt.annotate(txt, (x_b[i], Btotal_time[i]), xytext=(0, 5), textcoords='offset points', ha='center')
plt.xlabel('クライアントの数')
plt.ylabel('YCSB-Bにおける\n一命令あたりのレイテンシー(Millsec/命令)')
plt.legend()
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/latencyB/latency_plot-B.png')
plt.close()

# Data points for YCSB-C
x_c = [1, 2, 4, 8, 16, 32, 64]
Cread_time = [1.146807, 1.170556, 1.195731, 1.250501, 1.362784, 1.721079, 2.638301]
write_latency_ns = [0, 0, 0, 0, 0, 0, 0]  # すべてのクライアント数で0

# Creating the plot for YCSB-C
plt.figure(figsize=(10, 5))
plt.plot(x_c, Cread_time, marker='o', label='latency')
for i, txt in enumerate(x_c):
    plt.annotate(txt, (x_c[i], Cread_time[i]), xytext=(0, 5), textcoords='offset points', ha='center')
plt.xlabel('クライアントの数')
plt.ylabel('YCSB-Cにおける\n一命令あたりのレイテンシー(Millsec/命令)')
plt.legend()
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/latencyC/latency_plot-C.png')
plt.close()