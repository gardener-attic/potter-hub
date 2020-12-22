import * as React from "react";
import { IFooterConfig } from "shared/Config";
import logo from "../../img/logo.png";

interface IFooterProps {
  config?: IFooterConfig
}

const Footer: React.SFC<IFooterProps> = props => {
  const sectionOne = props.config?.sectionOne
  const sectionTwo = props.config?.sectionTwo
  return (
    <footer className="osFooter bg-dark type-color-reverse-anchor-reset">
      <div className="container padding-h-big padding-v-bigger">
        <div className="row collapse-b-phone-land">
          <div className="col-8">
            <h4 className="inverse margin-reset">
              <img src={logo} alt="Logo" title="Logo" className="osFooter__logo" />
            </h4>
          </div>
          <div className="col-2">
            {sectionTwo && (
              <>
                <p className="type-medium margin-reset">{sectionTwo.title}</p>
                {sectionTwo.links.map(item => (
                  <p className="type-small margin-reset">
                    <a href={item.href} target="_blank" rel="noopener noreferrer">
                      {item.title}
                    </a>
                  </p>
                ))}
              </>
            )}
          </div>
          <div className="col-2">
            {sectionOne && (
              <>
                <p className="type-medium margin-reset">{sectionOne.title}</p>
                {sectionOne.links.map(item => (
                  <p className="type-small margin-reset">
                    <a href={item.href} target="_blank" rel="noopener noreferrer">
                      {item.title}
                    </a>
                  </p>
                ))}
              </>
            )}
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
