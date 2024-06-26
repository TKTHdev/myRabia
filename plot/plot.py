import matplotlib.pyplot as plt

# Data points for YCSB-C
x_c = [1, 2, 3,  8, 16, 24, 32, 48, 64]
y_c = [57618/30, 106487/30, 177518/30, 385978/30, 571756/30, 604704/30 ,636320/30, 849903/30, 1006773/30]

# Creating the plot for YCSB-C
plt.figure(figsize=(10, 5))
plt.plot(x_c, y_c, marker='o')  # Line plot with circle markers
plt.title('Throughput vs. Number of Clients [YCSB-C]')
plt.xlabel('Number of Clients')
plt.ylabel('Throughput')
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/throughput_plot-C.png')
plt.close()  # Close the figure to prevent interference with the next plot

# Data points for YCSB-B
x_b = [1, 2, 3, 4, 8, 16, 24, 32, 48, 64]
y_b = [29298/30, 46016/30, 62446/30, 66546/30,65821/30, 66792/30, 69188/30, 68087/30, 68666/30, 69389/30]

# Creating the plot for YCSB-B
plt.figure(figsize=(10, 5))
plt.plot(x_b, y_b, marker='o')  # Line plot with circle markers
plt.title('Throughput vs. Number of Clients [YCSB-B]')
plt.xlabel('Number of Clients')
plt.ylabel('Throughput')
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/throughput_plot-B.png')
plt.close()  # Close the figure to ensure no memory issues or overlapsimport matplotlib.pyplot as plt



x_a = [1, 2, 3, 4, 8, 16, 24, 32, 48, 64]
y_a = [5884/30, 6553/30, 6549/30, 6969/30, 6961/30, 6971/30, 6784/30, 6833/30, 6791/30, 6909/30]


#Creating the plot for YCSB-A
plt.figure(figsize=(10, 5))
plt.plot(x_a, y_a, marker='o')  # Line plot with circle markers
plt.title('Throughput vs. Number of Clients [YCSB-A]')
plt.xlabel('Number of Clients')
plt.ylabel('Throughput')
plt.savefig('/home/Jamiroq/Documents/GitHub/myRabia/plot/throughput_plot-A.png')
plt.close()  # Close the figure to ensure no memory issues or overlapsimport matplotlib.pyplot as plt
