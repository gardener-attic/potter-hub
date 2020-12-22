import * as React from "react";

import "./CardGrid.css";

export interface ICardGridProps {
  className?: string;
  children?: React.ReactChildren | React.ReactNode | string;
}

interface ICardGridState {
  scrollEndIndex: number
}

const SCROLL_THRESHOLD_TOP_PIXEL = 120;
const SCROLL_THRESHOLD_BOTTOM_TO_LOAD = 0.8;
const ITEMS_PER_SCROLL = 32

class CardGrid extends React.Component<ICardGridProps, ICardGridState> {
  constructor(props: any) {
    super(props);
    this.state = {
      scrollEndIndex: ITEMS_PER_SCROLL,
    };
  }

  public componentDidMount() {
    window.addEventListener("scroll", this.onScroll);
  }

  public componentWillUnmount() {
    window.removeEventListener("scroll", this.onScroll)
  }

  public render() {
    return (
      <div className={`CardGrid padding-v-big ${this.props.className || ""}`}>
        {React.Children.toArray(this.props.children).slice(0, this.state.scrollEndIndex)}
      </div>
    );
  }

  private onScroll = (): void => {
    // calculates, if the current scroll position further than the SCROLL_THRESHOLD_BOTTOM_TO_LOAD threshold
    const isScrollAtBottom = window.innerHeight + document.documentElement.scrollTop >
      document.documentElement.offsetHeight * SCROLL_THRESHOLD_BOTTOM_TO_LOAD;

    // calculates, if the current scroll position is above the SCROLL_THRESHOLD_TOP_PIXEL
    const isScrollAtTop = document.documentElement.scrollTop < SCROLL_THRESHOLD_TOP_PIXEL;

    if (isScrollAtBottom && this.state.scrollEndIndex < React.Children.count(this.props.children)) {
      this.setState((state) => ({
        scrollEndIndex: state.scrollEndIndex + ITEMS_PER_SCROLL
      }))
    }

    if (isScrollAtTop && this.state.scrollEndIndex > ITEMS_PER_SCROLL) {
      this.setState((state) => ({
        scrollEndIndex: ITEMS_PER_SCROLL
      }))
    }
  }
}

export default CardGrid;
