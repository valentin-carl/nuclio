import os
from datetime import datetime

import matplotlib.pyplot as plt
import pandas as pd
import seaborn as sns
import numpy as np
import matplotlib.dates as mdates
from matplotlib.ticker import MaxNLocator
from enum import Enum
import matplotlib.lines as mlines


class Modes(Enum):
    NORMAL = 'normal'
    ASYNC = 'async'


class Scheduler(Enum):
    DEADLINE = 'deadline'
    IDLE = 'idle'
    BULK = 'bulk'


class ChartFactory:
    def __init__(self):
        self.fig, self.ax = plt.subplots()


        self.custom_colors = {
            Scheduler.DEADLINE.value: "#15F5BA",
            Scheduler.IDLE.value: "#836FFF",
            Scheduler.BULK.value: "#E26EE5",
            Modes.ASYNC.value: "#D16BA5",
            f"{Modes.ASYNC.value}_cpu": "#D16BA5",
            f"{Modes.ASYNC.value}_requests": "#FF8F80",
            Modes.NORMAL.value: "#86A8E7",
            f"{Modes.NORMAL.value}_cpu": "#86A8E7",
            f"{Modes.NORMAL.value}_requests": "#00DBED",

            "cpu_usage": "#FF9671",
        }
        self.interval = '50s'

        self.modes = ['async', 'normal']

    def finalize_and_save_chart(self, filename, title = ""):
        if title:
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
        self.finalize_and_save_chart(filename, title)

    def create_histogram(self, data, title, filename):
        df = pd.DataFrame(list(data.items()), columns=['scheduler_type', 'Count'])
        palette = [self.custom_colors.get(x) for x in df['scheduler_type'].unique()]
        sns.barplot(x='scheduler_type', y='Count', data=df, palette=palette)
        self.finalize_and_save_chart(filename, title)

    def create_line_chart(self, data, title, filename, include_scheduler=True):
        fig, ax1 = plt.subplots()

        if include_scheduler:
            deadline_data = data[data['scheduler_type'] == 'deadline']
            resampled_deadline_df = deadline_data.set_index('timestamp').resample('5s').size().reset_index(name='requests')
            resampled_deadline_df['scheduler_type'] = 'deadline'  # Add the scheduler type column back

            non_deadline_data = data[data['scheduler_type'] != 'deadline']
            resampled_non_deadline_df = non_deadline_data.set_index('timestamp').groupby('scheduler_type').resample(self.interval).size().reset_index(name='requests')

            resampled_df = pd.concat([resampled_deadline_df, resampled_non_deadline_df], ignore_index=True)
            sns.lineplot(x='timestamp', y='requests', hue='scheduler_type', data=resampled_df, ax=ax1,
                         palette=self.custom_colors)
            ax1.set_ylim(0, 100)
        else:
            resampled_df = data.set_index('timestamp').resample(self.interval).size().reset_index(name='requests')
            sns.lineplot(x='timestamp', y='requests', data=resampled_df, ax=ax1, color=self.custom_colors['normal'])

        plt.legend(title='Legend', loc='upper left')
        ax1.set_xlabel('Timestamp')
        ax1.set_ylabel('Requests per Minute')
        ax1.tick_params(axis='y')
        ax1.xaxis.set_major_formatter(mdates.DateFormatter('%H:%M'))
        ax1.xaxis.set_major_locator(MaxNLocator(6))

        data['cpu_usage'] = pd.to_numeric(data['cpu_usage'], errors='coerce')
        data = data.set_index('timestamp')

        resampled_cpu_df = data['cpu_usage'].resample('min').mean().reset_index()

        ax2 = ax1.twinx()
        sns.lineplot(x='timestamp', y='cpu_usage', label='CPU Usage', data=resampled_cpu_df, ax=ax2, color=self.custom_colors['cpu_usage'])
        ax2.set_ylabel('CPU Usage (%)')
        ax2.tick_params(axis='y')
        ax2.set_ylim(0, 100)
        ax1.set_xlim(resampled_df['timestamp'].min(), resampled_df['timestamp'].max())

        plt.xticks(rotation=180)

        self.finalize_and_save_chart(filename, title)

    def create_timeline_chart(self, data, title, filename):
        avg_durations = data.groupby('scheduler_type').agg({
            'async_duration': 'mean',
            'sync_duration': 'mean',
            'exec_duration': 'mean'
        }).reset_index()

        avg_durations_melted = avg_durations.melt(id_vars='scheduler_type',
                                                  value_vars=['async_duration', 'sync_duration', 'exec_duration'],
                                                  var_name='Phase', value_name='Average Duration')

        sns.barplot(x='Phase', y='Average Duration', hue='scheduler_type', data=avg_durations_melted, palette=self.custom_colors)
        plt.xticks(rotation=45)
        self.finalize_and_save_chart(filename, title)

    def create_cpu_histogram(self, async_data, normal_data, filename):
        async_avg_cpu = async_data['cpu_usage'].mean()
        normal_avg_cpu = normal_data['cpu_usage'].mean()
        df = pd.DataFrame({'Mode': self.modes, 'Average CPU Usage': [async_avg_cpu, normal_avg_cpu]})

        sns.barplot(x='Mode', y='Average CPU Usage', data=df, palette=self.custom_colors)
        self.finalize_and_save_chart(filename)

    def create_average_deviation_histogram(self, async_data, normal_data, filename):
        async_avg_req = async_data.groupby(pd.Grouper(key='timestamp', freq=self.interval)).size().mean()
        normal_avg_req = normal_data.groupby(pd.Grouper(key='timestamp', freq=self.interval)).size().mean()

        async_deviation = (async_data.groupby(pd.Grouper(key='timestamp', freq=self.interval)).size() - async_avg_req).abs().mean()
        normal_deviation = (normal_data.groupby(pd.Grouper(key='timestamp', freq=self.interval)).size() - normal_avg_req).abs().mean()

        df = pd.DataFrame({
            'Mode': self.modes,
            'Average Absolute Deviation': [async_deviation, normal_deviation]
        })

        sns.barplot(x='Mode', y='Average Absolute Deviation', data=df, palette=self.custom_colors)
        self.finalize_and_save_chart(filename)

    def create_execution_time_chart(self, async_data, normal_data, filename):
        async_avg_exec = async_data['exec_duration'].mean()
        normal_avg_exec = normal_data['exec_duration'].mean()

        df = pd.DataFrame({
            'Mode': self.modes,
            'Average Execution Time': [async_avg_exec, normal_avg_exec]
        })

        sns.barplot(x='Mode', y='Average Execution Time', data=df, palette=self.custom_colors)
        plt.xticks(rotation=45)
        self.finalize_and_save_chart(filename)

    def create_requests_per_minute_chart(self, async_data, normal_data, filename):
        async_avg_req = async_data.groupby(pd.Grouper(key='timestamp', freq='1Min')).size().mean()
        normal_avg_req = normal_data.groupby(pd.Grouper(key='timestamp', freq='1Min')).size().mean()

        df = pd.DataFrame({
            'Mode': self.modes,
            'Average Requests per Minute': [async_avg_req, normal_avg_req]
        })

        sns.barplot(x='Mode', y='Average Requests per Minute', data=df, palette=self.custom_colors)
        plt.xticks(rotation=45)
        self.finalize_and_save_chart(filename)

    def create_line_chart_comparison(self, async_data, normal_data, title, filename):
        fig, ax1 = plt.subplots()

        async_resampled = async_data.set_index('timestamp').resample(self.interval).agg({'cpu_usage': 'mean'})
        async_resampled['requests'] = async_data.set_index('timestamp').resample(self.interval).size()
        async_resampled['Mode'] = Modes.ASYNC.value
        async_resampled = async_resampled.reset_index(drop=True)

        normal_resampled = normal_data.set_index('timestamp').resample(self.interval).agg({'cpu_usage': 'mean'})
        normal_resampled['requests'] = normal_data.set_index('timestamp').resample(self.interval).size()
        normal_resampled['Mode'] = Modes.NORMAL.value
        normal_resampled = normal_resampled.reset_index(drop=True)

        combined_data = pd.concat([async_resampled, normal_resampled]).reset_index()

        cpu_palette = {Modes.ASYNC.value: self.custom_colors[f"{Modes.ASYNC.value}_cpu"],
                       Modes.NORMAL.value: self.custom_colors[f"{Modes.NORMAL.value}_cpu"]}
        request_palette = {Modes.ASYNC.value: self.custom_colors[f"{Modes.ASYNC.value}_requests"],
                           Modes.NORMAL.value: self.custom_colors[f"{Modes.NORMAL.value}_requests"]}

        sns.lineplot(x='index', y='requests', hue='Mode', data=combined_data, ax=ax1, palette=request_palette)

        ax2 = ax1.twinx()
        sns.lineplot(x='index', y='cpu_usage', hue='Mode', data=combined_data, ax=ax2, palette=cpu_palette, legend=False)
        ax2.set_ylabel('CPU Usage (%)')

        legend_elements = [
            mlines.Line2D([], [], color=self.custom_colors[f'{Modes.ASYNC.value}_requests'], label='Async Req.'),
            mlines.Line2D([], [], color=self.custom_colors[f'{Modes.NORMAL.value}_requests'], label='Normal Req.'),
            mlines.Line2D([], [], color=self.custom_colors[Modes.ASYNC.value], label='Async CPU'),
            mlines.Line2D([], [], color=self.custom_colors[Modes.NORMAL.value], label='Normal CPU')
        ]
        ax1.legend(handles=legend_elements, loc='upper right', title='Legend')

        ax1.set_xlabel('Relative Time Index')
        ax1.set_ylabel('Requests per Interval')

        plt.title(title)
        plt.xticks(rotation=45)
        fig.tight_layout()
        ax2.set_ylim(0, normal_resampled['requests'].max())
        ax1.set_xlim(0, 18)
        ax1.xaxis.set_major_locator(MaxNLocator(5))

        self.finalize_and_save_chart(filename, title)


def read_log_file(file_path):
    with open(file_path, 'r') as file:
        data = file.readlines()
    return data


def get_scheduler_count(dataframe):
    scheduler_counts = {}
    for scheduler_type in dataframe['scheduler_type']:
        scheduler_counts[scheduler_type] = scheduler_counts.get(scheduler_type, 0) + 1
    return scheduler_counts


def parse_log_data_for_line_chart(log_data):
    data = []
    for line in log_data:
        parts = line.strip().split(' - ')
        timestamp_str = parts[0]
        function_name = parts[1]
        scheduler_type = parts[2]
        async_incoming = pd.to_datetime(parts[3])
        sync_processing = pd.to_datetime(parts[4])
        exec_start = pd.to_datetime(parts[5])
        exec_stop = pd.to_datetime(parts[6])

        timestamp = datetime.strptime(timestamp_str, '%Y/%m/%d %H:%M:%S')

        data.append({
            'timestamp': timestamp,
            'function_name': function_name,
            'scheduler_type': scheduler_type,
            'async_incoming': async_incoming,
            'sync_processing': sync_processing,
            'exec_start': exec_start,
            'exec_stop': exec_stop,
            'async_duration': (sync_processing - async_incoming).total_seconds(),
            'sync_duration': (exec_start - sync_processing).total_seconds(),
            'exec_duration': (exec_stop - exec_start).total_seconds()
        })

    return pd.DataFrame(data)


def parse_normal_log_data(log_data):
    data = []
    for line in log_data:
        parts = line.strip().split(' - ')
        timestamp_str = parts[0]
        function_name = parts[1]
        exec_start = pd.to_datetime(parts[2])
        exec_stop = pd.to_datetime(parts[3])

        timestamp = datetime.strptime(timestamp_str, '%Y/%m/%d %H:%M:%S')
        data.append({
            'timestamp': timestamp,
            'function_name': function_name,
            'exec_start': exec_start,
            'exec_stop': exec_stop,
            'exec_duration': (exec_stop - exec_start).total_seconds()
        })

    return pd.DataFrame(data)


def merge_cpu_data(dataframe, cpu_log_data):
    cpu_data = []
    for line in cpu_log_data:
        parts = line.strip().split(' ')
        timestamp_str = parts[0] + ' ' + parts[1]
        cpu_usage = float(parts[4].strip('[]%'))
        timestamp = datetime.strptime(timestamp_str, '%Y/%m/%d %H:%M:%S')

        cpu_data.append({
            'timestamp': timestamp,
            'cpu_usage': cpu_usage
        })

    cpu_df = pd.DataFrame(cpu_data)

    merged_df = pd.merge_asof(dataframe.sort_values('timestamp'),
                              cpu_df.sort_values('timestamp'),
                              on='timestamp',
                              direction='nearest')
    return merged_df

chart_factory = ChartFactory()

general_log_data = read_log_file('./logs/async.log')
cpu_log_data = read_log_file('./logs/cpu-usage.log')
normal_log_data = read_log_file('./logs/normal.log')


async_dataframe = parse_log_data_for_line_chart(general_log_data)
async_dataframe = merge_cpu_data(async_dataframe, cpu_log_data)
normal_dataframe = parse_normal_log_data(normal_log_data)
normal_dataframe = merge_cpu_data(normal_dataframe, cpu_log_data)

log_counts = get_scheduler_count(async_dataframe)

chart_factory.create_chart('pie', log_counts, 'Log Type Distribution', 'async/pie_chart.png')
chart_factory.create_chart('histogram', log_counts, 'Log Type Distribution', 'async/histogram_chart.png')
chart_factory.create_chart('line', async_dataframe, 'Activity Monitor', 'async/line_chart.png')

chart_factory.create_timeline_chart(async_dataframe, 'Phase Durations by Scheduler', 'async/timeline_chart')
chart_factory.create_line_chart(normal_dataframe, 'Activity', 'normal/line_chart.png', include_scheduler=False)
chart_factory.create_line_chart_comparison(async_dataframe, normal_dataframe, 'Requests per Second', 'comparison/line_chart.png')

chart_factory.create_cpu_histogram(async_dataframe, normal_dataframe, 'comparison/cpu_histogram.png')
chart_factory.create_execution_time_chart(async_dataframe, normal_dataframe,  'comparison/execution_time_chart.png')
chart_factory.create_requests_per_minute_chart(async_dataframe, normal_dataframe,  'comparison/requests_per_minute_chart.png')
chart_factory.create_average_deviation_histogram(async_dataframe, normal_dataframe, 'comparison/requests_deviation_histogram.png')
