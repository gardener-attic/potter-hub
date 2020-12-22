import { getRandomName } from "./namegenerator";

it("generates name in the correct format", () => {
    const name = getRandomName()
    expect(name).toMatch(/\b[a-z]+-[a-z]+\b/g);
});