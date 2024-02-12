import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
from matplotlib.lines import Line2D

# load `./data/measurements.csv` into a DataFrame
measurements = pd.read_csv('./data/measurements.csv')
measurements = measurements.rename(columns={'latency': 'latency_sand'})

# load `./data/measurements-vanilla.csv` into a DataFrame
measurements_vanilla = pd.read_csv('./data/measurements-vanilla.csv')
measurements_vanilla = measurements_vanilla.rename(columns={'latency': 'latency_nuclio'})

# merge into a single DataFrame using the value column as the key
merged = pd.merge(measurements, measurements_vanilla, on='value')
#print(merged.head())

sns.set_style('whitegrid')
sns.set_context('notebook')

sns.scatterplot(x='value', y='latency_sand', hue="latency_sand" , data=merged, palette=["#7ad151"])
sns.scatterplot(x='value', y='latency_nuclio', hue="latency_nuclio", data=merged, palette=["#414487"])

plt.xlabel("Number of Function Calls in Request")
plt.ylabel("End-to-End Latency (ms)")

legend_labels = ['Sand', 'Nuclio']
legend_handles = [Line2D([0], [0], marker='o', color='w', markerfacecolor='#7ad151', markersize=10),
                  Line2D([0], [0], marker='o', color='w', markerfacecolor='#414487', markersize=10)]
plt.legend(legend_handles, legend_labels, title="Legend", loc="upper left", title_fontsize="12", fontsize="10", facecolor='white', framealpha=1)

#plt.show()
plt.savefig('latency.pdf')
