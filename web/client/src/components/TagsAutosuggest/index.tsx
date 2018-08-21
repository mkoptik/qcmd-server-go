import * as React from "react";
import Autosuggest = require("react-autosuggest");
import {ChangeEvent, InputProps, SuggestionSelectedEventData} from "react-autosuggest";
import "./style.css";

interface TagsAutosuggestProps {
    tags: string[][]
}

interface TagsAutosuggestState {
    inputValue: string,
    selectedSuggestion: string[]
}

export class TagsAutosuggest extends React.Component<TagsAutosuggestProps, TagsAutosuggestState> {

    constructor(props) {
        super(props);
        this.state = {
            inputValue: "",
            selectedSuggestion: null
        };
    }

    componentDidMount() {
        // TODO: FETCH ALL TAGS
    }

    render() {

        const inputProps: InputProps<string[]> = {
            value: this.state.inputValue,
            onChange: this.inputChange,
            placeholder: "Search executable"
        };

        return <Autosuggest
            suggestions={this.props.tags}
            onSuggestionsClearRequested={() => {}}
            onSuggestionsFetchRequested={() => {}}
            getSuggestionValue={suggestion => suggestion.join(" > ")}
            renderSuggestion={this.renderTagSuggestion}
            inputProps={inputProps}
            onSuggestionSelected={this.suggestionSelected} />
    }

    renderTagSuggestion = (suggestion: string[]) => suggestion.join(" > ");

    suggestionSelected = (e, data: SuggestionSelectedEventData<string[]>) => {
        this.setState({
            selectedSuggestion: data.suggestion
        })
    };

    inputChange = (event: React.FormEvent<any>, params?: ChangeEvent) => {
        console.log(params);
        this.setState({
            inputValue: params.newValue
        })
    }

}