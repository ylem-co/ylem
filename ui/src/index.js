import 'react-app-polyfill/ie11'; // For IE 11 support
import 'react-app-polyfill/stable';
//import './polyfill'
import React from 'react';
import App from './App';
import reportWebVitals from './reportWebVitals';
import * as serviceWorker from './serviceWorker';
import log from 'loglevel';

import { createRoot } from 'react-dom/client';

//import { icons } from './assets/icons'

import { Provider } from 'react-redux'
import store from './store'

//React.icons = icons

log.setLevel("trace", false)
log.info("Starting up...")

const container = document.getElementById('root');
const root = createRoot(container); // createRoot(container!) if you use TypeScript
root.render(<Provider store={store}><App/></Provider>);

//ReactDOM.render(
//  <Provider store={store}>
//    <App/>
//  </Provider>, 
//  document.getElementById('root')
//);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: http://bit.ly/CRA-PWA
serviceWorker.unregister();

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
