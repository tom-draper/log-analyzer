type Data = {
    extraction: {
        params: LineParams[];
        failed: string[];
    };
    config: {
        tokens: string[];
        patterns: string[];
    };
};

type LineParams = {
    [token: string]: string;
};

type TokenValueFreq = {
    [token: string]: ValueFreq;
};

type ValueFreq = {
    [value: string]: number
}