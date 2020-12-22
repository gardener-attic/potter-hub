import { isArray } from "lodash";
import * as React from "react";
import { IBasicFormParam, IBasicFormSliderParam } from "shared/types";
import TextParam from "./TextParam";

import { getValue } from "../../../shared/schema";
import "./BasicDeploymentForm.css";
import BooleanParam from "./BooleanParam";
import SliderParam from "./SliderParam";
import Subsection from "./Subsection";

export interface IBasicDeploymentFormProps {
  params: IBasicFormParam[];
  handleBasicFormParamChange: (
    p: IBasicFormParam,
  ) => (e: React.FormEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => void;
  handleValuesChange: (value: string) => void;
  appValues: string;
}

class BasicDeploymentForm extends React.Component<IBasicDeploymentFormProps> {
  public render() {
    return (
      <div className="margin-t-normal">
        {this.props.params.map((param, i) => {
          const id = `${param.path}-${i}`;
          return (
            <div key={id}>
              {this.renderParam(param, id, this.props.handleBasicFormParamChange)}
              <hr />
            </div>
          );
        })}
      </div>
    );
  }

  private isHidden = (param: IBasicFormParam) => {
    const hidden = param.hidden;
    switch (typeof hidden) {
      case "string":
        // If hidden is a string, it points to the value that should be true
        return getValue(this.props.appValues, hidden) === true;
      case "object":
        // Two type of supported objects
        // A single condition: {value: string, path: any}
        // An array of conditions: {conditions: Array<{value: string, path: any}, operator: string}
        if (hidden.conditions?.length > 0) {
          // If hidden is an object, a different logic should be applied
          // based on the operator
          switch (hidden.operator) {
            case "and":
              // Every value matches the referenced
              // value (via jsonpath) in all the conditions
              return hidden.conditions.every(
                c => getValue(this.props.appValues, c.path) === c.value,
              );
            case "or":
              // It is enough if the value matches the referenced
              // value (via jsonpath) in any of the conditions
              return hidden.conditions.some(
                c => getValue(this.props.appValues, c.path) === c.value,
              );
            case "nor":
              // Every value mismatches the referenced
              // value (via jsonpath) in any of the conditions
              return hidden.conditions.every(
                c => getValue(this.props.appValues, c.path) !== c.value,
              );
            default:
              // we consider 'and' as the default operator
              return hidden.conditions.every(
                c => getValue(this.props.appValues, c.path) === c.value,
              );
          }
        } else {
          return getValue(this.props.appValues, hidden.path) === hidden.value;
        }
      case "undefined":
        return false;
    }
  };

  private renderParam = (
    param: IBasicFormParam,
    id: string,
    handleBasicFormParamChange: (
      p: IBasicFormParam,
    ) => (e: React.FormEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => void,
  ) => {
    let paramComponent: JSX.Element = <></>;
    // If the type of the param is an array, represent it as its first type
    const type = isArray(param.type) ? param.type[0] : param.type;
    switch (type) {
      case "boolean":
        paramComponent = (
          <BooleanParam
            label={param.title || param.path}
            handleBasicFormParamChange={handleBasicFormParamChange}
            id={id}
            param={param}
          />
        );
        break;
      case "object": {
        paramComponent = (
          <Subsection
            label={param.title || param.path}
            handleValuesChange={this.props.handleValuesChange}
            appValues={this.props.appValues}
            renderParam={this.renderParam}
            param={param}
          />
        );
        break;
      }
      case "string": {
        if (param.render === "slider") {
          const p = param as IBasicFormSliderParam;
          paramComponent = (
            <SliderParam
              label={param.title || param.path}
              handleBasicFormParamChange={handleBasicFormParamChange}
              id={id}
              param={param}
              min={p.sliderMin || 1}
              max={p.sliderMax || 1000}
              unit={p.sliderUnit || ""}
            />
          );
          break;
        }
        if (param.render === "textArea") {
          paramComponent = (
            <TextParam
              label={param.title || param.path}
              handleBasicFormParamChange={handleBasicFormParamChange}
              id={id}
              param={param}
              inputType="textarea"
            />
          );
          break;
        }
      }
      default:
        paramComponent = (
          <TextParam
            label={param.title || param.path}
            handleBasicFormParamChange={handleBasicFormParamChange}
            id={id}
            param={param}
            inputType={type === "integer" ? "number" : "string"}
          />
        );
    }
    return (
      <div key={id} hidden={this.isHidden(param)}>
        {paramComponent}
      </div>
    );
  };
}

export default BasicDeploymentForm;
