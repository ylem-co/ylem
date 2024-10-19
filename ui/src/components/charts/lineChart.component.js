import React from 'react';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  LineElement,
  Title,
  Tooltip,
  Legend,
  PointElement
} from 'chart.js';
import { Line } from 'react-chartjs-2';

ChartJS.register(
  CategoryScale,
  LinearScale,
  LineElement,
  Title,
  Tooltip,
  Legend,
  PointElement
);

const isDarkThemeEnabled = localStorage.getItem('darkTheme') !== "false";

export const lineDataSample = {
  labels: [],
  datasets: [
    {
      label: ' Average duration',
      data: [],
      borderColor: isDarkThemeEnabled ? 'rgb(206, 150, 250)' : 'rgb(176, 96, 239)',
      color: isDarkThemeEnabled ? 'rgb(206, 150, 250)' : 'rgb(176, 96, 239)',
      pointRadius: 6,
      pointHoverRadius: 10,
    },
  ],
};

export const lineValueDataSample = {
  labels: [],
  datasets: [
    {
      label: ' Value',
      data: [],
      borderColor: isDarkThemeEnabled ? 'rgb(206, 150, 250)' : 'rgb(176, 96, 239)',
      color: isDarkThemeEnabled ? 'rgb(206, 150, 250)' : 'rgb(176, 96, 239)',
      pointRadius: 6,
      pointHoverRadius: 10,
    },
  ],
};

export const lineChartOptions = {
  scales: {
    y: {
      stacked: true,
      ticks: {
        beginAtZero: true,
        backdropColor: isDarkThemeEnabled ? '#333b43' : '#fff',
        color: isDarkThemeEnabled ? 'white' : 'black',
      },
      grid:{
        drawBorder:false,
        color: isDarkThemeEnabled ? 'rgb(172, 175, 177, 0.2)' : 'rgb(172, 175, 177, 0.4)',
      },
    },
    x: {
      stacked: true,
      ticks: {
        backdropColor: isDarkThemeEnabled ? '#333b43' : '#fff',
        color: isDarkThemeEnabled ? 'white' : 'black',
      },
      grid:{
        drawBorder:false,
        color: isDarkThemeEnabled ? 'rgb(172, 175, 177, 0.2)' : 'rgb(172, 175, 177, 0.4)',
      },
    },
  },
  maintainAspectRatio: false,
  plugins: {
    tooltip: {
      callbacks: {

      }
    },
    legend: {
      display: false
    }
  }
};

export const lineChartValueOptions = {
  scales: {
    y: {
      stacked: true,
      ticks: {
        beginAtZero: true,
        backdropColor: isDarkThemeEnabled ? '#333b43' : '#fff',
        color: isDarkThemeEnabled ? 'white' : 'black',
      },
      grid:{
        drawBorder:false,
        color: isDarkThemeEnabled ? 'rgb(172, 175, 177, 0.2)' : 'rgb(172, 175, 177, 0.4)',
      },
    },
    x: {
      stacked: true,
      ticks: {
        display: false,
        backdropColor: isDarkThemeEnabled ? '#333b43' : '#fff',
        color: isDarkThemeEnabled ? 'white' : 'black',
      },
      grid:{
        drawBorder:false,
        color: isDarkThemeEnabled ? 'rgb(172, 175, 177, 0.2)' : 'rgb(172, 175, 177, 0.4)',
      },
    },
  },
  maintainAspectRatio: false,
  plugins: {
    tooltip: {
      callbacks: {

      }
    },
    legend: {
      display: false
    }
  }
};

const LineChart = ({ data, options }) => (
  <>
    <Line data={data} options={options} />
  </>
);

export default LineChart;