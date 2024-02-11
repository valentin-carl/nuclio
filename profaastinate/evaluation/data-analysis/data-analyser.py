import os
from datetime import datetime

import matplotlib.pyplot as plt
import pandas as pd
import seaborn as sns


class ChartFactory:
    def __init__(self):
        self.fig, self.ax = plt.subplots()

        self.custom_colors = {"deadline": "#D6BCFA",
                              "bulk": "#93C5FD"}

    def finalize_and_save_chart(self, title, filename):
        plt.title(title)
        filename = f"./charts/{filename}"
        if os.path.isfile(filename):
            os.remove(filename)
        plt.savefig(filename, bbox_inches='tight')
        plt.close()

    def create_chart(self, chart_type, data, title, filename):
        if chart_type == 'pie':
            self.create_pie_chart(data, title, filename)
        elif chart_type == 'histogram':
            self.create_histogram(data, title, filename)
        elif chart_type == 'line':
            self.create_line_chart(data, title, filename)
        else:
            raise ValueError(f"Unknown chart type: {chart_type}")

    def create_pie_chart(self, data, title, filename):
        labels = data.keys()
        sizes = data.values()
        colors = [self.custom_colors.get(x, "#333333") for x in labels]
        self.ax.pie(sizes, labels=labels, autopct='%1.1f%%', startangle=90, colors=colors)
        self.ax.axis('equal')
        self.finalize_and_save_chart(title, filename)

    def create_histogram(self, data, title, filename):
        df = pd.DataFrame(list(data.items()), columns=['scheduler_type', 'Count'])
        palette = [self.custom_colors.get(x, "#333333") for x in df['scheduler_type'].unique()]
        sns.barplot(x='scheduler_type', y='Count', data=df, palette=palette)
        self.finalize_and_save_chart(title, filename)

    def create_line_chart(self, data, title, filename):
        # Grouping by 'scheduler_type', resampling by second, and counting occurrences
        resampled_df = data.set_index('timestamp').groupby('scheduler_type').resample('s').size().reset_index(
            name='requests')

        palette = [self.custom_colors.get(x, "#333333") for x in resampled_df['scheduler_type'].unique()]
        sns.lineplot(x='timestamp', y='requests', hue='scheduler_type', data=resampled_df,palette=palette)
        plt.xticks(rotation=45)
        self.finalize_and_save_chart(title, filename)


def read_log_file(file_path):
    with open(file_path, 'r') as file:
        data = file.readlines()
    return data


def parse_log_data(log_data):
    log_counts = {}
    for line in log_data:
        log_type = line.strip().split(' ')[-1]
        if log_type in log_counts:
            log_counts[log_type] += 1
        else:
            log_counts[log_type] = 1
    return log_counts


def parse_log_data_for_line_chart(log_data):
    data = []
    for line in log_data:
        parts = line.strip().split(' ')
        timestamp_str = ' '.join(parts[:2])
        scheduler_type = parts[-1]

        timestamp = datetime.strptime(timestamp_str, '%Y/%m/%d %H:%M:%S')
        data.append({'timestamp': timestamp, 'scheduler_type': scheduler_type})

    return pd.DataFrame(data)


log_data = read_log_file('./logs/async.log')
log_counts = parse_log_data(log_data)

chart_factory = ChartFactory()
chart_factory.create_chart('pie', log_counts, 'Log Type Distribution', 'pie_chart.png')
chart_factory.create_chart('histogram', log_counts, 'Log Type Distribution', 'histogram_chart.png')

log_data_for_line_chart = parse_log_data_for_line_chart(log_data)
chart_factory.create_chart('line', log_data_for_line_chart, 'Requests per Second by Scheduler', 'line_chart.png')
