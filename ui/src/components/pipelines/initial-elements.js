export default [
  {
    id: '1',
    type: 'query',
    data: {
      name: 'DB query',
    },
    position: { x: -200, y: 100 },
  },
  {
    id: '2',
    type: 'condition',
    data: {
      name: 'IF A > B',
    },
    position: { x: -90, y: 100 },
  },
  {
    id: '3',
    type: 'aggregator',
    data: {
      name: 'AVG(amount)',
    },
    position: { x: 50, y: 100 },
  },
  {
    id: '4',
    type: 'api_call',
    data: {
      name: 'Call API',
    },
    position: { x: 250, y: 100 },
  },
  {
    id: '5',
    type: 'transformer',
    data: {
      name: 'Data transformation',
    },
    position: { x: 50, y: 200 },
  },
  {
    id: '6',
    type: 'notification',
    data: {
      name: 'Send SMS',
    },
    position: { x: 250, y: 200 },
  },
  { id: 'e1-2', source: '2', target: '1' },
  { id: 'e2-3', source: '3', target: '2', label: 'true' },
  { id: 'e3-4', source: '4', target: '3' },
  { id: 'e2-5', source: '5', target: '2', label: 'false' },
  { id: 'e5-6', source: '6', target: '5' },
];
