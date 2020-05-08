import React from 'react';
import {Card, Row, Col} from "@canonical/react-components";

function DetailsCard(props) {
    return (
        <Card>
            {props.fields.map(f => {
                return (
                    <Row>
                        <Col size={2}>{f.label}:</Col>
                        <Col size={10} className="field_value">{f.value}</Col>
                    </Row>
                )
            })}
        </Card>
    );
}

export default DetailsCard;