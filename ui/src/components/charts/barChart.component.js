import React from 'react';

import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import { Bar } from 'react-chartjs-2';

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend
);

const isDarkThemeEnabled = localStorage.getItem('darkTheme') !== "false";

export const barDataSample = {
  labels: [],
  datasets: [
    {
      label: 'Successful runs',
      data: [],
      backgroundColor: 'rgba(2, 140, 217, 0.65)',
      color: isDarkThemeEnabled ? 'white' : 'black',
    },
    {
      label: 'Failures',
      data: [],
      backgroundColor: 'rgba(254, 113, 136, 0.65)',
      color: isDarkThemeEnabled ? 'white' : 'black',
    },
  ],
};

const options = {
  plugins: {
    legend: {
      labels: {
        color: isDarkThemeEnabled ? 'white' : 'black',
      }, 
    },
  },
  scales: {
    y: {
      stacked: true,
      ticks: {
        beginAtZero: true,
        backdropColor: isDarkThemeEnabled ? '#333b43' : '#fff',
        color: isDarkThemeEnabled ? 'white' : 'black',
      }
    },
    x: {
      stacked: true,
      ticks: {
        backdropColor: isDarkThemeEnabled ? '#333b43' : '#fff',
        color: isDarkThemeEnabled ? 'white' : 'black',
      }
    }
  },
  maintainAspectRatio: false,
};

const BarChart = ({ data }) => (
  <>
    <Bar data={data} options={options} />
  </>
);

export default BarChart;