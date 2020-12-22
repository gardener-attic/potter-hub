import * as React from "react";
import { AlertOctagon, Info } from "react-feather";
import Modal from "react-modal";
import CreatableSelect from "react-select/creatable";

import { INamespaceState } from "../../reducers/namespace";
import { definedNamespaces } from "../../shared/Namespace";
import LoadingWrapper from "../LoadingWrapper";

import "./NamespaceSelector.css";

interface INamespaceSelectorProps {
  namespace: INamespaceState;
  defaultNamespace: string;
  onChange: (ns: string) => any;
  fetchNamespaces: () => void;
  createNamespace: (ns: string) => void;
  clearNamespaceError: () => void;
}

interface INamespaceSelectorState {
  showModal: boolean;
  newNamespace: string;
  createNamespaceRequestWasSent: boolean;
}

const customSelectStyle = {
  input: (provided: any) => ({
    ...provided,
    minHeight: "1px",
    maxHeight: "36px",
    paddingTop: "0",
    paddingBottom: "0",
    margin: "0",
  })
};

class NamespaceSelector extends React.Component<INamespaceSelectorProps, INamespaceSelectorState> {
  public state: INamespaceSelectorState = {
    showModal: false,
    newNamespace: "",
    createNamespaceRequestWasSent: false
  };

  public componentDidMount() {
    this.props.fetchNamespaces();
  }

  public render() {
    const {
      namespace: { current, namespaces },
      defaultNamespace,
    } = this.props;
    const options =
      namespaces.length > 0
        ? namespaces.map(n => ({ value: n, label: n }))
        : [{ value: defaultNamespace, label: defaultNamespace }];
    const allOption = { value: definedNamespaces.all, label: "All Namespaces" };
    options.unshift(allOption);
    const selected = current || defaultNamespace;
    const value =
      selected === definedNamespaces.all ? allOption : { value: selected, label: selected };

    return (
      <div className="NamespaceSelector margin-r-normal">
        <label className="NamespaceSelector__label type-tiny">NAMESPACE</label>
        {this.renderNewNamespaceModal()}
        <CreatableSelect
          className="NamespaceSelector__select type-small"
          classNamePrefix="NamespaceSelector"
          value={value}
          options={options}
          multi={false}
          onChange={this.handleNamespaceChange}
          formatCreateLabel={this.promptTextCreator}
          clearable={false}
          onCreateOption={this.handleNewOptionClick}
          styles={customSelectStyle}
        />
      </div>
    );
  }

  public renderNewNamespaceModal = () => {
    const {
      namespace: { errorMsg, isFetching },
    } = this.props;
    const {
      newNamespace,
      showModal,
    } = this.state;

    let modal: JSX.Element
    if (isFetching) {
      modal = <Modal
        className="centered-modal CreateNamespaceModal"
        isOpen={showModal}
        onRequestClose={this.resetModalState}
        contentLabel="Modal"
        shouldCloseOnEsc={false}
        shouldCloseOnOverlayClick={false}
      >
        <div className="row confirm-dialog-loading-info">
          <div className="col-8 loading-legend">Loading, please wait</div>
          <div className="col-4">
            <LoadingWrapper />
          </div>
        </div>
      </Modal>
    }
    else if (errorMsg) {
      modal = <Modal
        className="centered-modal CreateNamespaceModal"
        isOpen={showModal}
        onRequestClose={this.resetModalState}
        contentLabel="Modal"
        shouldCloseOnEsc={true}
        shouldCloseOnOverlayClick={true}
        shouldFocusAfterRender={true}
      >
        <div>
          <div className="alert alert-error AlertBox" role="alert">
            <div className="AlertBoxIconContainer">
              <AlertOctagon />
            </div>
            <span>Cannot create namespace: {errorMsg}</span>
          </div>
          <div className="margin-t-normal CreateNamespaceErrorDialogBtnContainer">
            <button className="button" onClick={this.resetModalState}>
              Ok
          </button>
          </div>
        </div>
      </Modal>
    } else if (this.state.createNamespaceRequestWasSent) {
      modal = <Modal
        className="centered-modal CreateNamespaceModal"
        isOpen={showModal}
        onRequestClose={this.handleSuccesModalClose}
        contentLabel="Modal"
        shouldCloseOnEsc={true}
        shouldCloseOnOverlayClick={true}
      >
        <div>
          <div className="alert alert-success AlertBox" role="alert">
            <div className="AlertBoxIconContainer">
              <Info />
            </div>
            <span>Namespace "{newNamespace}" was succesfully created</span>
          </div>
          <div className="margin-t-normal CreateNamespaceErrorDialogBtnContainer">
            <button className="button" onClick={this.handleSuccesModalClose}>
              Ok
        </button>
          </div>
        </div>
      </Modal>
    }
    else {
      modal = <Modal
        className="centered-modal CreateNamespaceModal"
        isOpen={showModal}
        onRequestClose={this.resetModalState}
        contentLabel="Modal"
        shouldCloseOnEsc={true}
        shouldCloseOnOverlayClick={true}
      >
        <div>
          <div className="margin-b-normal">
            Create namespace "{newNamespace}"?
        </div>
          <div className="margin-t-normal button-row">
            <button className="button" onClick={this.resetModalState}>
              Cancel
          </button>
            <button
              className="button button-primary button-danger"
              type="submit"
              onClick={this.createNamespace}
            >
              Ok
          </button>
          </div>
        </div>
      </Modal>
    }
    return modal
  }

  public handleSuccesModalClose = () => {
    this.props.onChange(this.state.newNamespace)
    this.resetModalState()
  }

  public resetModalState = () => {
    this.props.clearNamespaceError()
    this.setState({
      showModal: false,
      createNamespaceRequestWasSent: false,
    });
  }

  public handleNewOptionClick = (option: any) => {
    this.setState({
      showModal: true,
      newNamespace: option,
    });
  }

  public createNamespace = () => {
    this.props.createNamespace(this.state.newNamespace)
    this.setState({
      createNamespaceRequestWasSent: true
    });
  };

  private handleNamespaceChange = (value: any) => {
    if (value) {
      this.props.onChange(value.value);
    }
  };

  private promptTextCreator = (text: string) => {
    return `Create namespace "${text}"`;
  };

}

export default NamespaceSelector;
