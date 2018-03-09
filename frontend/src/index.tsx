import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './app';
import registerServiceWorker from './registerServiceWorker';
import { createStore } from 'redux';
import myApp from './reducers'

let store = createStore(myApp);

window.onload = function render() {
  ReactDOM.render(
    <div className='container'>
      <App store={store} />
    </div>,
    document.getElementById('root')
  );
}

// store.subscribe(render);

// render();
// ReactDOM.render(<App />, document.getElementById('root'));
// registerServiceWorker();
