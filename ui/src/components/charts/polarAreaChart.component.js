import React from 'react';
import {
  Chart as ChartJS,
  RadialLinearScale,
  ArcElement,
  Tooltip,
  Legend,
  PointElement
} from 'chart.js';
import { PolarArea } from 'react-chartjs-2';

ChartJS.register(
  RadialLinearScale,
  ArcElement,
  Tooltip,
  Legend,
  PointElement
);

const isDarkThemeEnabled = localStorage.getItem('darkTheme') !== "false";

export const polarAreaDashboardData = {
  datasets: [
    {
      label: 'Pipelines',
      backgroundColor: [
        'rgba(2, 140, 217, 0.65)',
        'rgba(247, 186, 27, 0.65)',
        'rgba(254, 113, 136, 0.65)',
        'rgba(66, 189, 111, 0.65)',
        'rgba(131, 81, 245, 0.65)',
        'rgba(144, 191, 92, 0.65)',
        'rgba(186, 78, 182, 0.65)',
      ],
      borderWidth: 0,
    },
  ],
};

const options = {
    plugins: {
      legend: {
        display: false,
      },
    },
    scales: {
      r: {
        grid: {
          color: isDarkThemeEnabled ? 'rgb(172, 175, 177, 0.2)' : 'rgb(172, 175, 177, 0.4)',
        },
        ticks: {
          backdropColor: isDarkThemeEnabled ? '#333b43' : '#fff',
          color: isDarkThemeEnabled ? 'white' : 'black',
        }
      },
    }
  };

const PolarAreaChart = ({ data }) => (
  <>
    <PolarArea 
      data={data} 
      options={options}
    />
  </>
);

export default PolarAreaChart;
