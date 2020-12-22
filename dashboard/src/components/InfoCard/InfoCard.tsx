import * as React from "react";

import { Link } from "react-router-dom";
import placeholder from "../../img/placeholder.png";
import triangle from "../../img/triangle.png";
import Card, { CardContent, CardIcon } from "../Card";
import "./InfoCard.css";

export interface IServiceInstanceCardProps {
  title: string;
  info: string | JSX.Element;
  link?: string;
  icon?: string;
  banner?: string;
  description?: string | JSX.Element;
  tag1Class?: string;
  tag1Content?: string | JSX.Element;
  tag2Class?: string;
  tag2Content?: string | JSX.Element;
  tag3Content?: string | JSX.Element;
  highlight?: boolean;
  tagContainerCSS?: React.CSSProperties | undefined
}

const InfoCard: React.SFC<IServiceInstanceCardProps> = props => {
  const {
    title,
    link,
    info,
    description,
    tag1Content,
    tag1Class,
    tag2Content,
    tag2Class,
    tag3Content,
    banner,
    highlight,
    tagContainerCSS,
  } = props;
  const icon = props.icon ? props.icon : placeholder;

  let tag3: JSX.Element | undefined;

  if (typeof tag3Content === "string") {
    let parsedMetadata;
    try {
      parsedMetadata = JSON.parse(tag3Content);
    } catch {
      // do nothing
    }
    if (parsedMetadata?.bomName) {
      tag3 = (
        <Link to={`/clusterboms/${parsedMetadata.bomName}`}>
          <span
            className={
              "ListItem__content__info_tag type-small type-color-white padding-t-tiny padding-h-normal ListItem__content__info_tag-3"
            }
          >
            BoM: {parsedMetadata.bomName}
          </span>
        </Link>
      );
    } else if (tag3Content.toLocaleLowerCase() === "controllermanaged") {
      tag3 = (
        <span
          className={
            "ListItem__content__info_tag type-small type-color-white padding-t-tiny padding-h-normal ListItem__content__info_tag-3"
          }
        >
          managed by BoM
        </span>
      );
    }
  } else if (tag3Content) {
    tag3 = tag3Content;
  }

  return (
    <Card responsive={true} className="ListItem">
      <Link to={link || "#"} title={title} className="ListItem__header">
        <CardIcon icon={icon} />
        {banner && <p className="ListItem__banner">{banner}</p>}
        {highlight && (
          <img
            src={triangle}
            className="ListItem__triangle"
            title="Featured Chart"
            alt="Featured Chart"
          />
        )}
      </Link>
      <CardContent>
        <div className="ListItem__content">
          <Link to={link || "#"} title={title}>
            <h3 className="ListItem__content__title type-big">{title}</h3>
          </Link>
          {description}
          <div className="ListItem__content__info">
            <p className="margin-reset type-small padding-t-tiny type-color-light-blue">{info}</p>
            <div style={tagContainerCSS}>
              {tag1Content && (
                <span
                  className={`ListItem__content__info_tag ListItem__content__info_tag-1 type-small type-color-white padding-t-tiny padding-h-normal ${tag1Class ||
                    ""}`}
                >
                  {tag1Content}
                </span>
              )}
              {tag2Content && (
                <span
                  className={`ListItem__content__info_tag ListItem__content__info_tag-2 type-small type-color-white padding-t-tiny padding-h-normal ${tag2Class ||
                    ""}`}
                >
                  {tag2Content}
                </span>
              )}
              {tag3}
            </div>
          </div>
        </div>
      </CardContent>
      {props.children}
    </Card>
  );
};

export default InfoCard;
