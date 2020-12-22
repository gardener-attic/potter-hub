import * as React from "react";
import { IHelpMenuConfig } from "shared/Config";

import "./HelpMenu.css";

interface IHelpMenuProps {
  appVersion: string;
  appName: string;
  config?: IHelpMenuConfig;
  controllerAppVersion: string;
}

interface IHelpMenuState {
  isOpen: boolean;
}

class HelpMenu extends React.Component<IHelpMenuProps, IHelpMenuState> {
  public state: IHelpMenuState = {
    isOpen: false,
  };

  private timeOutId: number | undefined;
  private helpMenuRef: React.RefObject<HTMLElement>;
  private helpMenuButtonRef: React.RefObject<HTMLButtonElement>;

  constructor(props: IHelpMenuProps) {
    super(props);
    this.helpMenuRef = React.createRef();
    this.helpMenuButtonRef = React.createRef();
  }

  public render() {
    return (
      <React.Fragment>
        <button
          className="icon-button"
          ref={this.helpMenuButtonRef}
          onClick={this.onClickHandler}
          aria-haspopup="true"
          aria-expanded={this.state.isOpen}
        >
          <svg xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 0 24 24" width="24">
            <path d="M0 0h24v24H0z" fill="none" />
            <path
              className="icon-button-icon"
              d="M11 18h2v-2h-2v2zm1-16C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm0-14c-2.21 0-4 1.79-4 4h2c0-1.1.9-2 2-2s2 .9 2 2c0 2-3 1.75-3 5h2c0-2.25 3-2.5 3-5 0-2.21-1.79-4-4-4z"
            />
          </svg>
        </button>
        {this.state.isOpen && (
          <section
            className="help-menu elevation-2"
            tabIndex={-1}
            ref={this.helpMenuRef}
            onBlur={this.onBlurHandler}
            onFocus={this.onFocusHandler}
          >
            <div className="header-container padding-big">
              <p className="heading">{this.props.appName}</p>
              <p className="version">Hub Version: {this.props.appVersion}</p>
              <p className="version">Controller Version: {this.props.controllerAppVersion}</p>
            </div>
            {this.props.config?.links && (
              this.props.config?.links.map(item => (
                <div className="entry padding-normal">
                  <a href={item.href} target="_blank" rel="noopener noreferrer">
                    <div className="link">
                      {item.title}
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        height="24"
                        viewBox="0 0 24 24"
                        width="24"
                      >
                        <path d="M0 0h24v24H0z" fill="none" />
                        <path
                          className="link-icon"
                          d="M19 19H5V5h7V3H5c-1.11 0-2 .9-2 2v14c0 1.1.89 2 2 2h14c1.1 0 2-.9 2-2v-7h-2v7zM14 3v2h3.59l-9.83 9.83 1.41 1.41L19 6.41V10h2V3h-7z"
                        />
                      </svg>
                    </div>
                  </a>
                </div>
              ))
            )}
          </section>
        )}
      </React.Fragment>
    );
  }

  public componentDidUpdate() {
    if (this.state.isOpen) {
      if (this.helpMenuRef.current != null) {
        this.helpMenuRef.current.focus();
      }
    }
  }

  private onClickHandler = (e: React.MouseEvent<HTMLButtonElement>) => {
    this.setState(prevState => {
      return { isOpen: !prevState.isOpen };
    });
  };

  private onBlurHandler = (e: React.FocusEvent) => {
    if (e.relatedTarget !== this.helpMenuButtonRef.current) {
      this.timeOutId = setTimeout(() => {
        this.setState({
          isOpen: false,
        });
      });
    }
  };

  private onFocusHandler = () => {
    clearTimeout(this.timeOutId);
  };
}

export default HelpMenu;
