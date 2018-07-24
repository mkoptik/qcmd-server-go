import * as React from 'react';
import "./style.css";
import * as qs from "querystring"
import {Command} from "../../models/Command";
import Axios from "axios";
import {RouteComponentProps} from "react-router";
import {TagsAutosuggest} from "../TagsAutosuggest";

interface AppProps extends RouteComponentProps<any> {
    initialCommands: Command[],
    allTags: string[][],
    apiUrl: string
}

interface AppState {
    search: string,
    foundCommands: Command[],
}

export class App extends React.Component<AppProps, AppState> {

    constructor(props) {
        super(props);
        const queryString = qs.parse(props.location.search.replace("?", ""));
        this.state = {
            search: queryString["search"] as string || "",
            foundCommands: []
        };
    }

    componentWillReceiveProps(newProps: AppProps) {
        const bla = qs.parse(newProps.location.search);
        console.log(bla)
    }

    inputTextChanged = (value) => {
        this.setState({ search: value }, () => {
            Axios.get<Command[]>(`${this.props.apiUrl}/search`, { params: { search: value} })
                .then(response => { this.setState({ foundCommands: response.data }) })
        });
    };

    render() {
        const commands = this.state.foundCommands.length > 0 ? this.state.foundCommands : this.props.initialCommands;
        return (
            <div>
                <div className="inputs-wrap">
                    <div className="tag-input"><TagsAutosuggest tags={this.props.allTags} /></div>
                    <div className="search-input">
                        <input type="text" onChange={e => this.inputTextChanged(e.currentTarget.value)} placeholder="Search command" value={this.state.search} />
                    </div>
                </div>
                <div className="commands-list">
                    {commands.map(c => this.renderCommand(c))}
                </div>
            </div>
        );
    }

    renderCommand = (command: Command) => <div key={command.label} className="command">
        {command.label}
        <div className="command-text">{command.commandText}</div>
    </div>;


}
