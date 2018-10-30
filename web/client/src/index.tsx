import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { App } from './components/App';
import {Route} from "react-router";
import {BrowserRouter} from "react-router-dom";

ReactDOM.render(
    <BrowserRouter>
        <Route path="/" render={p => <App initialCommands={[]} allTags={[]} apiUrl={"http://localhost:8888"} {...p} />} />
    </BrowserRouter>
    ,
    document.getElementById("root")
);
