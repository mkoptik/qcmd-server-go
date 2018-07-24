import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { App } from './components/App';
import {Command} from "./models/Command";
import {Route} from "react-router";
import {BrowserRouter} from "react-router-dom";

export function init(rootElementId: string, apiUrl: string, initialCommands: Command[], allTags: string[][]) {
    ReactDOM.render(
        <BrowserRouter>
            <Route path="/" render={p => <App initialCommands={initialCommands} allTags={allTags} apiUrl={apiUrl} {...p} />} />
        </BrowserRouter>
        ,
        document.getElementById(rootElementId)
    );
}
