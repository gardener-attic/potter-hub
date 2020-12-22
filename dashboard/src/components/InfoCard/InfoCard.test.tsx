import { shallow } from "enzyme";
import * as React from "react";
import { Link } from "react-router-dom";

import InfoCard from "./InfoCard";

it("should render a Card", () => {
  const wrapper = shallow(
    <InfoCard
      title="foo"
      info="foobar"
      link="/a/link/somewhere"
      icon="an-icon.png"
      tag1Class="blue"
      tag1Content="database"
      tag2Class="red"
      tag2Content="running"
      tag3Content="controllermanaged"
    />,
  );
  expect(wrapper).toMatchSnapshot();
});

it("should generate a dummy link if it's not provided", () => {
  const wrapper = shallow(
    <InfoCard
      title="foo"
      info="foobar"
      icon="an-icon.png"
      tag1Class="blue"
      tag1Content="database"
      tag2Class="red"
      tag2Content="running"
      tag3Content="controllermanaged"
    />,
  );

  const links = wrapper.find(Link);
  expect(links.length).toBe(2);
  links.forEach(l => expect(l.props()).toMatchObject({ to: "#" }));
});

it("should avoid tags if they are not defined", () => {
  const wrapper = shallow(
    <InfoCard title="foo" info="foobar" link="/a/link/somewhere" icon="an-icon.png" />,
  );
  expect(wrapper.find(".ListItem__content__info_tag")).not.toExist();
});

it("should parse JSX elements in the tags", () => {
  const tag1 = <span className="tag1">tag1</span>;
  const tag2 = <span className="tag2">tag2</span>;
  const wrapper = shallow(
    <InfoCard
      title="foo"
      info="foobar"
      link="/a/link/somewhere"
      icon="an-icon.png"
      tag1Class="blue"
      tag1Content={tag1}
      tag2Class="red"
      tag2Content={tag2}
      tag3Content="controllermanaged"
    />,
  );
  const tag1Elem = wrapper.find(".tag1");
  expect(tag1Elem).toExist();
  expect(tag1Elem.text()).toBe("tag1");
  const tag2Elem = wrapper.find(".tag2");
  expect(tag2Elem).toExist();
  expect(tag2Elem.text()).toBe("tag2");
});

it("should parse a description as text", () => {
  const wrapper = shallow(<InfoCard title="foo" info="foobar" description="a description" />);
  expect(wrapper.find(".ListItem__content").text()).toContain("a description");
});

it("should parse a description as JSX.Element", () => {
  const desc = <div className="description">This is a description</div>;
  const wrapper = shallow(<InfoCard title="foo" info="foobar" description={desc} />);
  expect(wrapper.find(".description").text()).toBe("This is a description");
});

it("should render a banner if exists", () => {
  const wrapper = shallow(<InfoCard title="foo" info="foobar" banner="this is important!" />);
  const banner = wrapper.find(".ListItem__banner");
  expect(banner).toExist();
  expect(banner.text()).toBe("this is important!");
});

it("should render a highlight", () => {
  const wrapper = shallow(<InfoCard title="foo" info="foobar" highlight={true} />);
  const highlight = wrapper.find(".ListItem__triangle");
  expect(highlight).toExist();
});

it("should not render a highlight", () => {
  const wrapper = shallow(<InfoCard title="foo" info="foobar" />);
  const highlight = wrapper.find(".ListItem__triangle");
  expect(highlight).not.toExist();
});

it("should render a link to the respective clusterbom if tag3 contains appropriate metadata", () => {
  const metadata = {
    bomName: "testbom",
  };

  const wrapper = shallow(
    <InfoCard title="foo" info="foobar" tag3Content={JSON.stringify(metadata)} />,
  );

  const tag3Elem = wrapper.find(".ListItem__content__info_tag-3");
  expect(tag3Elem).toExist();
  expect(tag3Elem.text()).toBe("BoM: testbom");
  expect(tag3Elem.parent().type()).toBe(Link);
  expect(tag3Elem.parent().props().to).toEqual("/clusterboms/testbom");
});
