type Data = {
    extraction: Extraction[];
    config: Config;
    types: DataTypes
    failed: Failed
    locations: Locations
};

type Extraction = {
    params: Params;
    line: string
    lineNumber: number
    pattern: string
}

type Config = {
    tokens: string[];
    patterns: string[];
    dependencies?: string[];
    conversions?: {
        [token: string]: {
            token: string,
            multiplier: number
        }
    };
}

type DataTypes = {
    [token: string]: {
        [type: string]: number
    }
}

type Failed = {
    [lineNumber: number]: string
}

type Locations = {
    [ipAddress: string]: string
}

type Params = {
    [token: string]: string | number;
};

type TokenValueFreq = {
    [token: string]: ValueCount;
};

type ValueCount = {
    [value: string]: number
}

type ValueCounts = {
    [value: string]: number[]
}